use std::{
    collections::{HashMap, HashSet},
    path::{Path, PathBuf},
    process::Command,
};

use anyhow::{bail, Context};
use clap::{Parser, Subcommand};
use dialoguer::{theme::ColorfulTheme, Input, Select};
use regex::Regex;
use serde::{Deserialize, Serialize};

use crate::collection::{Collection, EncodeMetadata};

mod collection;
mod metadata;

// TODO(patrik):
//   - Feat: Check collection for corruption
//   - Feat: Repair corruption
//   - Feat: Import should detect changes
//   - TODO: When reencoding for mobile check the bitrate
//   - Feat: User confirmation on create-config before writing the config

// const TARGET_MOBILE_BIT_RATE: usize = 192000;

// File Structure:
//   - artists.json
//   - {SOME_FOLDER_NAME}
//     - album.json : Album Metadata
//     - images/
//     - tracks/
//       - full/
//          - track00.flac
//          - track01.flac
//       - mobile/
//          - track00.mp3
//          - track01.mp3
// artists.json:
//   - ID
//   - Name
//   - Picture
// album.json:
//   - Name
//   - Artist ID (null if multiple artists)
//   - CoverArt
//   - Tracks:
//       - Track Number
//       - Name
//       - full: track00.flac
//       - mobile: track00.mp3
//

#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
    #[command(subcommand)]
    command: ArgCommand,
}

#[derive(Subcommand, Debug)]
enum ArgCommand {
    List { collection: PathBuf },
    Import { collection: PathBuf, path: PathBuf },
    CreateConfig { path: Option<PathBuf> },
}

#[derive(Deserialize, Debug)]
struct FormatMetadata {
    format_name: String,
    bit_rate: String,

    tags: Option<TagsMetadata>,
}

#[derive(Deserialize, Debug)]
struct TagsMetadata {
    #[serde(alias = "TITLE")]
    title: String,
    #[serde(alias = "ARTIST")]
    artist: String,
    #[serde(alias = "ALBUM")]
    album: String,
    track: Option<String>,
    disc: Option<String>,
    album_artist: Option<String>,
}

impl TagsMetadata {
    fn track(&self) -> Option<usize> {
        let track = if let Some((left, _right)) =
            self.track.as_ref()?.split_once("/")
        {
            left.parse::<usize>().ok()?
        } else {
            self.track.as_ref()?.parse::<usize>().ok()?
        };

        Some(track)
    }

    fn album(&self) -> String {
        return if let Some(disc) = &self.disc {
            if let Some((left, _right)) = disc.split_once("/") {
                if left == "1" {
                    self.album.clone()
                } else {
                    format!("{} (Disc {})", self.album, left)
                }
            } else {
                if disc == "1" {
                    self.album.clone()
                } else {
                    format!("{} (Disc {})", self.album, disc)
                }
            }
        } else {
            self.album.clone()
        };
    }
}

fn ffprobe<P>(path: P) -> anyhow::Result<FormatMetadata>
where
    P: AsRef<Path>,
{
    let path = path.as_ref();

    // fwav fprobe -print_format json -v quiet -show_format -show_streams test.mp3

    let output = Command::new("ffprobe")
        .arg("-print_format")
        .arg("json")
        .arg("-v")
        .arg("quiet")
        .arg("-show_format")
        .arg(path)
        .output()
        .context("Failed to run ffprobe 'Is ffprobe installed?'")?;

    if !output.status.success() {
        bail!("ffprobe error: {:?}", path);
    }

    let s = std::str::from_utf8(&output.stdout)
        .context("Failed to convert ffprobe stdout to str")?;

    let data = serde_json::from_str::<serde_json::Value>(s)
        .context("Failed to convert ffprobe data to json")?;
    let format = data
        .get("format")
        .context("No format object in ffprobe data")?;
    // let tags = format.get("tags").context("No tags object in format")?;

    let metadata = serde_json::from_value::<FormatMetadata>(format.clone())
        .context("Failed to convert format object to FormatMetadata")?;

    Ok(metadata)
}

#[derive(Serialize, Deserialize, Debug)]
struct TrackDef {
    num: usize,
    name: String,
    filename: String,
    artist: Option<String>,
}

#[derive(Serialize, Deserialize, Debug)]
struct DiscDef {
    disc_num: usize,
    dir: PathBuf,
    tracks: Vec<TrackDef>,
}

#[derive(Serialize, Deserialize, Debug)]
struct MultiDiscDef {
    name: String,
    artist: String,
    discs: Vec<DiscDef>,
}

#[derive(Serialize, Deserialize, Debug)]
struct NormalDef {
    name: String,
    artist: String,
    tracks: Vec<TrackDef>,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(tag = "type")]
enum Config {
    MultiDisc(MultiDiscDef),
    Normal(NormalDef),
}

fn list<P>(collection_path: P) -> anyhow::Result<()>
where
    P: AsRef<Path>,
{
    let collection = Collection::new(collection_path)?;
    for album in collection.albums() {
        let metadata = album.metadata();
        println!("{}: {}", metadata.id, metadata.name);
        for track in album.metadata().tracks.iter() {
            let artist = collection.artist_by_id(&track.artist_id).unwrap();
            println!(
                "  {:2} - {} - {}",
                track.track_num, artist.name, track.name
            );
        }
    }

    Ok(())
}

fn import<P>(collection_path: P, path: PathBuf) -> anyhow::Result<()>
where
    P: AsRef<Path>,
{
    let mut collection = Collection::new(collection_path)?;

    let mut album_toml_path = path.clone();
    album_toml_path.push("album.toml");

    if !album_toml_path.is_file() {
        return Err(anyhow::anyhow!("No 'album.toml' file found"));
    }

    let s = std::fs::read_to_string(album_toml_path)?;
    let def = toml::from_str::<Config>(&s)?;
    let Config::Normal(def) = def else { todo!() };
    println!("Def: {:#?}", def);

    let album_artist_index = collection.get_or_insert_artist(&def.artist);
    let album_artist = collection.artist_by_index(album_artist_index);
    // TODO(patrik): Remove clone?
    let album_artist = album_artist.id.clone();

    if let Some(album) = collection.add_album(&album_artist, &def.name) {
        println!("Album: {:#?}", album);
        for track in def.tracks {
            let mut path = path.clone();
            path.push(track.filename);

            let artist = track.artist.as_ref().unwrap_or(&def.artist);
            let artist = collection.get_or_insert_artist(&artist);

            let encode_metadata = EncodeMetadata {
                track: track.num,
                title: &track.name,
                album: &def.name,
                album_artist: album_artist_index,
                artist,
            };

            collection
                .process_track(album, path, encode_metadata)
                .unwrap();
        }
    } else {
        // What to do here?
        bail!("Album '{}' already exists", def.name);
    }

    // Check if the album exists
    // Process all the tracks

    collection.save_to_disk().context("Failed save to disk")?;

    Ok(())
}

enum SearchOption {
    Normal,
    MultiDisc,
}

fn get_all_tracks_from_dir<P>(
    path: P,
) -> anyhow::Result<Vec<(PathBuf, FormatMetadata)>>
where
    P: AsRef<Path>,
{
    let path = path.as_ref();
    let mut music_files = Vec::new();

    for entry in path.read_dir()? {
        let entry = entry?;
        let path = entry.path();

        if path.is_dir() {
            continue;
        }
        println!("Entry: {:?}", entry.path().display());

        if let Ok(metadata) = ffprobe(&path) {
            match metadata.format_name.as_str() {
                "flac" | "mov,mp4,m4a,3gp,3g2,mj2" | "mp3" | "wav" => {
                    music_files.push((path, metadata));
                }

                format => eprintln!("Unsupported format: {}", format),
            }
        }
    }

    Ok(music_files)
}

fn create_config_for_album(
    path: &PathBuf,
    gussed_album_name: &str,
    gussed_album_artist: &str,
) -> anyhow::Result<()> {
    struct FileTrack<'a> {
        path: &'a PathBuf,
        metadata: &'a FormatMetadata,
        file_title: String,
        file_artist: Option<String>,
    }

    let music_files = get_all_tracks_from_dir(&path)?;

    let mut tracks = HashMap::new();

    for (path, metadata) in music_files.iter() {
        let file_name = path.file_name().unwrap().to_str().unwrap();

        let (left, _) = file_name.rsplit_once(".").unwrap();
        let mut splits = left.split("-");
        let track_num = splits
            .next()
            .unwrap()
            .trim()
            .parse::<usize>()
            .with_context(|| {
                format!("File name missing track number: '{}'", file_name)
            })?;

        let file_title = splits.next().unwrap().trim().to_string();
        let file_artist = splits.next().map(|a| a.trim().to_string());

        tracks.insert(
            track_num,
            FileTrack {
                path,
                metadata,
                file_title,
                file_artist,
            },
        );
    }

    if tracks.len() <= 0 {
        bail!("No tracks found");
    }

    // TODO(patrik): Detect duplicated track titles

    let mut album_names = HashSet::new();
    album_names.insert(gussed_album_name);

    let mut album_artist_names = HashSet::new();
    album_artist_names.insert(gussed_album_artist);

    for (_, track) in tracks.iter() {
        if let Some(tags) = track.metadata.tags.as_ref() {
            album_names.insert(tags.album.as_str());
            if let Some(artist) = tags.album_artist.as_ref() {
                album_artist_names.insert(artist.as_str());
            }
        }
    }

    let album_name = if album_names.len() > 1 {
        let mut items = album_names.into_iter().collect::<Vec<_>>();
        items.sort();
        items.push("Custom Name");

        let selection = Select::with_theme(&ColorfulTheme::default())
            .with_prompt("Select Album Name")
            .items(&items)
            .default(0)
            .interact_opt();

        match selection? {
            Some(select) => {
                if select == items.len() - 1 {
                    let text: String =
                        Input::with_theme(&ColorfulTheme::default())
                            .with_prompt("Enter Custom Name")
                            .interact_text()?;
                    text
                } else {
                    items[select].to_string()
                }
            }
            None => {
                println!("Quitting...");
                return Ok(());
            }
        }
    } else if album_names.len() == 1 {
        album_names.into_iter().next().unwrap().to_string()
    } else {
        // TODO(patrik): Let the user choose
        println!("No metadata found, trying to guess");

        gussed_album_name.to_string()
    };

    let album_artist = if album_artist_names.len() > 1 {
        let mut items = album_artist_names.into_iter().collect::<Vec<_>>();
        items.sort();
        items.push("Various Artists");
        items.push("Custom Name");
        let selection = Select::with_theme(&ColorfulTheme::default())
            .with_prompt("Select Album Artist")
            .items(&items)
            .default(0)
            .interact_opt();

        match selection? {
            Some(select) => {
                if select == items.len() - 1 {
                    let text: String =
                        Input::with_theme(&ColorfulTheme::default())
                            .with_prompt("Enter Custom Name")
                            .interact_text()?;
                    text
                } else {
                    items[select].to_string()
                }
            }
            None => {
                println!("Quitting...");
                return Ok(());
            }
        }
    } else if album_artist_names.len() == 1 {
        album_artist_names.into_iter().next().unwrap().to_string()
    } else {
        // TODO(patrik): Let the user choose

        gussed_album_artist.to_string()
    };

    println!("Album Name: {}", album_name);
    println!("Album Artist: {}", album_artist);

    let mut tracks = tracks
        .into_iter()
        .map(|(track_num, track)| TrackDef {
            num: track_num,
            name: track
                .metadata
                .tags
                .as_ref()
                .map(|tags| tags.title.clone())
                .unwrap_or(track.file_title),
            filename: track
                .path
                .file_name()
                .unwrap()
                .to_str()
                .unwrap()
                .to_string(),
            artist: track
                .metadata
                .tags
                .as_ref()
                .map(|tags| tags.artist.clone())
                .or(track.file_artist)
                .or_else(|| Some(album_artist.clone())),
        })
        .collect::<Vec<_>>();

    tracks.sort_by(|l, r| l.num.cmp(&r.num));

    let config = Config::Normal(NormalDef {
        name: album_name,
        artist: album_artist.clone(),
        tracks,
    });

    let s = toml::to_string(&config).unwrap();
    println!("Config: {}", s);

    let mut output = PathBuf::from(path);
    output.push("album.toml");
    std::fs::write(output, s).unwrap();

    Ok(())
}

fn create_config(path: PathBuf, search: SearchOption) -> anyhow::Result<()> {
    match search {
        SearchOption::Normal => {
            let album_name = path.file_name().unwrap().to_str().unwrap();

            let artist_name = path
                .parent()
                .unwrap()
                .file_name()
                .unwrap()
                .to_str()
                .unwrap();

            create_config_for_album(&path, &album_name, artist_name)?;
        }

        SearchOption::MultiDisc => {
            let regex = Regex::new(r"(disc|cd)\s*(\d+)")?;

            for entry in path.read_dir()? {
                let entry = entry?;
                let entry = entry.path();

                if !entry.is_dir() {
                    continue;
                }

                let name = entry
                    .file_name()
                    .unwrap()
                    .to_str()
                    .unwrap()
                    .to_lowercase();
                if let Some(caps) = regex.captures(&name) {
                    let disk_num = caps[2].parse::<usize>()?;

                    let album_name =
                        path.file_name().unwrap().to_str().unwrap();

                    let album_name =
                        format!("{} (Disc {})", album_name, disk_num);

                    let artist_name = path
                        .parent()
                        .unwrap()
                        .file_name()
                        .unwrap()
                        .to_str()
                        .unwrap();

                    create_config_for_album(&entry, &album_name, artist_name)?;
                }
            }
        }
    }

    Ok(())
}

fn main() -> anyhow::Result<()> {
    ctrlc::set_handler(move || {
        let term = dialoguer::console::Term::stdout();
        let _ = term.show_cursor();
    })
    .context("Failed to set ctrl-c handler")?;

    let args = Args::parse();

    match args.command {
        ArgCommand::List { collection } => list(collection)?,
        ArgCommand::Import { collection, path } => import(collection, path)?,
        ArgCommand::CreateConfig { path } => {
            let options = ["Normal Album", "Multi Disc Album"];
            let selection = Select::with_theme(&ColorfulTheme::default())
                .with_prompt("Select type of album")
                .items(&options)
                .default(0)
                .interact()?;
            println!("Selected {}", options[selection]);

            let search_option = if selection == 1 {
                SearchOption::MultiDisc
            } else {
                SearchOption::Normal
            };

            create_config(
                path.unwrap_or(std::env::current_dir()?),
                search_option,
            )?
        }
    }

    Ok(())
}
