use std::{
    path::{Path, PathBuf},
    process::Command,
};

use anyhow::{bail, Context};
use clap::{Parser, Subcommand};
use serde::{Deserialize, Serialize};
use walkdir::DirEntry;

use crate::collection::{Collection, EncodeMetadata};

mod collection;

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
    collection: PathBuf,

    #[command(subcommand)]
    command: ArgCommand,
}

#[derive(Subcommand, Debug)]
enum ArgCommand {
    List {},
    Import { path: PathBuf },
}

fn is_hidden(entry: &DirEntry) -> bool {
    entry
        .file_name()
        .to_str()
        .map(|s| s.starts_with("."))
        .unwrap_or(false)
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
struct TrackMetadata {
    #[serde(rename = "trackNum")]
    track_num: usize,
    name: String,
    artist_id: String,

    quality_version: String,
    mobile_version: String,
}

#[derive(Serialize, Deserialize, Debug)]
struct AlbumMetadata {
    id: String,
    name: String,
    #[serde(rename = "artistId")]
    artist_id: Option<String>,
    #[serde(rename = "coverArt")]
    cover_art: String,

    tracks: Vec<TrackMetadata>,
}

#[derive(Serialize, Deserialize, Debug)]
struct ArtistMetadata {
    id: String,
    name: String,
    picture: String,
}

#[derive(Serialize, Deserialize, Debug)]
struct TrackDef {
    num: usize,
    name: String,
    filename: String,
    artist: Option<String>,
}

#[derive(Serialize, Deserialize, Debug)]
struct AlbumDef {
    name: String,
    artist: String,
    tracks: Vec<TrackDef>,
}

fn main() -> anyhow::Result<()> {
    let args = Args::parse();
    println!("Args: {:#?}", args);

    let mut collection = Collection::new(args.collection)?;

    match args.command {
        ArgCommand::List {} => {
            for album in collection.albums() {
                let metadata = album.metadata();
                println!("{}: {}", metadata.id, metadata.name);
                for track in album.metadata().tracks.iter() {
                    let artist =
                        collection.artist_by_id(&track.artist_id).unwrap();
                    println!(
                        "  {:2} - {} - {}",
                        track.track_num,
                        artist.name,
                        track.name
                    );
                }
            }
        }

        ArgCommand::Import { path } => {
            let mut def_toml_path = path.clone();
            def_toml_path.push("def.toml");

            if !def_toml_path.is_file() {
                return Err(anyhow::anyhow!("No 'def.toml' file found"));
            }

            let s = std::fs::read_to_string(def_toml_path)?;
            let def = toml::from_str::<AlbumDef>(&s)?;
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

            collection.save_to_disk().expect("Failed save to disk");
        }
    }

    // let combined_artist_name = "__COMBINED_ALBUMS__";
    // let combined_albums = [
    //     "We Will Rock You (Disc 1)",
    //     "We Will Rock You (Disc 2)",
    //     "Cyberpunk 2077",
    // ];
    //
    // let mut artists: HashMap<String, Artist> = HashMap::new();
    //
    // let mut insert = |artist: &str, metadata: &Metadata| {
    //     let artist =
    //         artists.entry(artist.to_string()).or_insert_with(|| Artist {
    //             name: artist.to_string(),
    //             albums: HashMap::new(),
    //         });
    //
    //     let albums =
    //         artist
    //             .albums
    //             .entry(metadata.album())
    //             .or_insert_with(|| Album {
    //                 name: metadata.album(),
    //                 songs: HashMap::new(),
    //             });
    //
    //     let track = metadata.track().unwrap_or(albums.songs.len());
    //     if albums.songs.contains_key(&track) {
    //         panic!("Mutliple songs with same tracks");
    //     }
    //     albums.songs.insert(
    //         track,
    //         Song {
    //             name: metadata.title.clone(),
    //         },
    //     );
    // };
    //
    // match args.command {
    //     ArgCommand::Scan { path } => {
    //         println!("Scanning {:?}", path);
    //
    //         for entry in WalkDir::new(path)
    //             .into_iter()
    //             .filter_entry(|e| !is_hidden(e))
    //             .filter_map(|e| e.ok())
    //         {
    //             if entry.path().is_dir() {
    //                 continue;
    //             }
    //
    //             if entry
    //                 .path()
    //                 .file_name()
    //                 .context("Failed to convert filename")?
    //                 .to_str()
    //                 .context("Failed to convert filename to str")?
    //                 == "tracks.txt"
    //             {
    //                 continue;
    //             }
    //
    //             // println!("Entry: {:?}", entry);
    //             match ffprobe(entry.path()) {
    //                 Ok(metadata) => {
    //                     if combined_albums.contains(&metadata.album.as_str()) {
    //                         insert(combined_artist_name, &metadata);
    //                     } else {
    //                         insert(&metadata.artist, &metadata);
    //                     }
    //                 }
    //                 Err(e) => eprintln!("Error: {:?}", e),
    //             }
    //         }
    //
    //         println!("{:#?}", artists);
    //     }
    // }

    Ok(())
}
