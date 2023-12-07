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
use walkdir::DirEntry;

use crate::collection::{Collection, EncodeMetadata};

mod collection;
mod metadata;

// TODO(patrik):
//   - Feat: Check collection for corruption
//   - Feat: Repair corruption
//   - Feat: Import should detect changes
//   - TODO: When reencoding for mobile check the bitrate

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

    let mut def_toml_path = path.clone();
    def_toml_path.push("def.toml");

    if !def_toml_path.is_file() {
        return Err(anyhow::anyhow!("No 'def.toml' file found"));
    }

    let s = std::fs::read_to_string(def_toml_path)?;
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

fn create_config(path: PathBuf) -> anyhow::Result<()> {
    let mut music_file = Vec::new();

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
                    music_file.push((path, metadata));
                }

                format => eprintln!("Unsupported format: {}", format),
            }
        }
    }

    struct FileTrack<'a> {
        path: &'a PathBuf,
        metadata: &'a FormatMetadata,
        file_title: String,
    }

    let mut tracks = HashMap::new();

    let reg = Regex::new(r"(\d+)[\s-]+(.*)\.\w+").unwrap();

    for (path, metadata) in music_file.iter() {
        let file_name = path.file_name().unwrap().to_str().unwrap();

        if let Some(caps) = reg.captures(file_name) {
            let track_num = caps[1].parse::<usize>()?;
            let name = caps[2].to_string();
            tracks.insert(
                track_num,
                FileTrack {
                    path,
                    metadata,
                    file_title: name,
                },
            );
        } else {
            println!("No track number in filename detected");
        }
    }

    if tracks.len() <= 0 {
        bail!("No tracks found");
    }

    // TODO(patrik): Detect duplicated track titles

    let mut album_names = HashSet::new();
    let mut album_artist_names = HashSet::new();

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
        let album_name =
            path.file_name().unwrap().to_str().unwrap().to_string();

        album_name
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
    } else if album_artist_names.len() == 1{
        album_artist_names.into_iter().next().unwrap().to_string()
    } else {
        // TODO(patrik): Let the user choose
        let artist_name = path
            .parent()
            .unwrap()
            .file_name()
            .unwrap()
            .to_str()
            .unwrap()
            .to_string();

        artist_name
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
    output.push("def.toml");
    std::fs::write(output, s).unwrap();

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
            create_config(path.unwrap_or(std::env::current_dir()?))?
        }
    }

    Ok(())
}
