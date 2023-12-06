use anyhow::{bail, Context, Result};
use std::{
    path::{Path, PathBuf},
    process::{Command, ExitStatus},
};

use crate::{
    ffprobe, AlbumMetadata, ArtistMetadata, EncodeMetadata, TrackMetadata,
};

#[derive(Debug)]
struct Artist {
    id: String,
    name: String,
}

#[derive(Debug)]
pub struct Track {
    track_num: usize,
    name: String,
    artist_index: usize,
}

impl Track {
    pub fn track_num(&self) -> usize {
        self.track_num
    }

    pub fn name(&self) -> &str {
        &self.name
    }

    pub fn artist(&self) -> usize {
        self.artist_index
    }
}

#[derive(Debug)]
pub struct Album {
    path: PathBuf,

    id: String,
    name: String,

    tracks: Vec<Track>,
}

impl Album {
    pub fn id(&self) -> &str {
        &self.id
    }

    pub fn name(&self) -> &str {
        &self.name
    }

    pub fn tracks(&self) -> &Vec<Track> {
        &self.tracks
    }
}

pub struct Collection {
    base: PathBuf,
    artists: Vec<Artist>,
    albums: Vec<Album>,
}

impl Collection {
    pub fn new<P>(path: P) -> Result<Self>
    where
        P: AsRef<Path>,
    {
        let path = path.as_ref().to_path_buf();

        let mut collection = Collection {
            base: path,
            artists: Vec::new(),
            albums: Vec::new(),
        };

        collection.load_artists()?;
        collection.load_albums()?;

        Ok(collection)
    }

    fn load_artists(&mut self) -> Result<()> {
        let mut artists_path = self.base.clone();
        artists_path.push("artists.json");

        if let Ok(s) = std::fs::read_to_string(artists_path) {
            let artists = serde_json::from_str::<Vec<ArtistMetadata>>(&s)
                .context("Failed to parse 'artists.json'")?;

            self.artists.try_reserve(artists.len())?;
            self.artists
                .extend(artists.into_iter().map(|artist| Artist {
                    id: artist.id,
                    name: artist.name,
                }));
        }

        Ok(())
    }

    fn load_albums(&mut self) -> Result<()> {
        let mut album_path = self.base.clone();
        album_path.push("albums");

        if !album_path.exists() {
            return Ok(());
        }

        for entry in album_path.read_dir()? {
            let entry = entry?;
            let entry = entry.path();
            println!("Entry: {:?}", entry);

            let mut album_json = entry.clone();
            album_json.push("album.json");

            let s = std::fs::read_to_string(album_json)?;
            let metadata = serde_json::from_str::<AlbumMetadata>(&s)?;
            println!("Metadata: {:#?}", metadata);

            let tracks: Vec<Track> = metadata
                    .tracks
                    .into_iter()
                    .map(|track| -> Result<Track> {
                        Ok(Track {
                            track_num: track.track_num,
                            name: track.name,
                            artist_index: self
                                .get_artist_by_id(&track.artist_id)
                                .with_context(|| {
                                    format!(
                                        "No artist with id: {} (collection might be corrupted)",
                                        track.artist_id
                                    )
                                })?,
                        })
                    })
                    .collect::<Result<_, _>>()?;

            self.albums.push(Album {
                path: entry,
                id: "".to_string(), //metadata.id,
                name: metadata.name,
                tracks,
            });
        }

        Ok(())
    }

    pub fn albums(&self) -> &Vec<Album> {
        &self.albums
    }

    // TODO(patrik): Check if artist already exists
    pub fn get_or_insert_artist(&mut self, name: &str) -> usize {
        let new_id = cuid2::CuidConstructor::new().with_length(18);

        let index = self.artists.len();
        self.artists.push(Artist {
            id: new_id.create_id(),
            name: name.to_string(),
        });

        index
    }

    fn get_artist_by_name(&self, name: &str) -> Option<&Artist> {
        None
    }

    fn get_artist_by_id(&self, id: &str) -> Option<usize> {
        self.artists
            .iter()
            .enumerate()
            .find(|(index, artist)| artist.id == id)
            .map(|(index, _)| index)
    }

    pub fn get_album_by_name(&self, name: &str) -> Option<()> {
        self.albums
            .iter()
            .find(|album| album.name == name)
            .map(|_| ())
    }

    fn get_album_by_name_mut(&mut self, name: &str) -> Option<&mut Album> {
        self.albums.iter_mut().find(|album| album.name == name)
    }

    fn generate_new_album(&self) -> Result<(String, PathBuf)> {
        let new_id = cuid2::CuidConstructor::new().with_length(18);

        let id = new_id.create_id();
        let mut path = self.base.clone();
        path.push("albums");
        path.push(&id);

        if path.exists() {
            // TODO(patrik): Just try again
            bail!("Album with id already exists");
        }

        std::fs::create_dir_all(&path)?;

        Ok((id, path))
    }

    pub fn add_album(&mut self, name: String) -> Option<usize> {
        if self.get_album_by_name(&name).is_some() {
            return None;
        }

        let (id, path) =
            self.generate_new_album().expect("Failed to generate album");

        println!("Id: {:?}", id);
        println!("Path: {:?}", path);

        let index = self.albums.len();
        self.albums.push(Album {
            path,
            id,
            name,
            tracks: Vec::new(),
        });

        Some(index)
    }

    pub fn process_track<P>(
        &mut self,
        album_index: usize,
        file_path: P,
        encode_metadata: EncodeMetadata,
    ) -> Result<()>
    where
        P: AsRef<Path>,
    {
        let file_path = file_path.as_ref();
        println!("Processing: {:?}", file_path);

        let metadata =
            ffprobe(file_path).context("Failed to process track")?;
        println!("Metadata: {:#?}", metadata);

        let album = &self.albums[album_index];
        let mut output = album.path.clone();
        output.push("tracks");

        std::fs::create_dir_all(&output)?;

        match metadata.format_name.as_str() {
            "wav" => {
                let bit_rate = metadata
                    .bit_rate
                    .parse::<usize>()
                    .context("Failed to parse bitrate")?;
                println!("Bitrate: {}", bit_rate);

                let mut out = output.clone();
                out.push("full");
                std::fs::create_dir_all(&out)?;
                out.push(format!("{}.flac", encode_metadata.track));

                let status =
                    self.encode_flac(file_path, out, &encode_metadata)?;
                println!("Flac Status: {:?}", status);

                let mut out = output.clone();
                out.push("mobile");
                std::fs::create_dir_all(&out)?;
                out.push(format!("{}.mp3", encode_metadata.track));

                let status =
                    self.encode_mp3(file_path, out, &encode_metadata)?;
                println!("Mp3 Status: {:?}", status);

                let album = &mut self.albums[album_index];
                album.tracks.push(Track {
                    track_num: encode_metadata.track,
                    name: encode_metadata.title.to_string(),
                    artist_index: 0,
                });

                Ok(())
            }
            name => Err(anyhow::anyhow!("Unsupported format '{}'", name)),
        }
    }

    // fn load_artists(&mut self) {}
    // fn save_artists(&mut self) {}
    //
    // fn load_albums(&mut self) {}
    // fn add_album(&mut self) {}

    pub fn save_to_disk(self) -> Result<()> {
        for album in self.albums {
            let metadata = AlbumMetadata {
                name: album.name,
                artist_id: None,
                cover_art: "".to_string(),
                tracks: album
                    .tracks
                    .into_iter()
                    .map(|track| TrackMetadata {
                        track_num: track.track_num,
                        name: track.name,
                        artist_id: self.artists[track.artist_index].id.clone(),
                        full: "".to_string(),
                        mobile: "".to_string(),
                    })
                    .collect(),
            };

            let mut metadata_path = album.path;
            metadata_path.push("album.json");

            let s = serde_json::to_string_pretty(&metadata)?;
            std::fs::write(metadata_path, s)?;
        }

        let artists = self.artists.into_iter().map(|artist| ArtistMetadata {
            id: artist.id,
            name: artist.name,
            picture: "".to_string(),
        }).collect::<Vec<_>>();

        let mut artist_metadata_path = self.base;
        artist_metadata_path.push("artists.json"); 

        let s = serde_json::to_string_pretty(&artists)?;
        std::fs::write(artist_metadata_path, s)?;

        Ok(())
    }

    fn encode_flac<I, O>(
        &self,
        input: I,
        output: O,
        metadata: &EncodeMetadata,
    ) -> Result<ExitStatus>
    where
        I: AsRef<Path>,
        O: AsRef<Path>,
    {
        let input = input.as_ref();
        let output = output.as_ref();

        let ext = output
            .extension()
            .context("Failed to get extention")?
            .to_str()
            .context("Failed to convert extention to str")?;
        assert_eq!(ext, "flac");

        let mut command = Command::new("ffmpeg");
        command.arg("-y").arg("-i").arg(input);

        command
            .arg("-metadata")
            .arg(format!("title={}", metadata.title));

        command
            .arg("-metadata")
            .arg(format!("album={}", metadata.album));

        command
            .arg("-metadata")
            .arg(format!("artist={}", self.artists[metadata.artist].name));

        command
            .arg("-metadata")
            .arg(format!("track={}", metadata.track));

        command.arg(output);

        let status = command
            .status()
            .context("Failed to run ffmpeg 'Is ffmpeg installed?'")?;

        Ok(status)
    }

    fn encode_mp3<I, O>(
        &self,
        input: I,
        output: O,
        metadata: &EncodeMetadata,
    ) -> Result<ExitStatus>
    where
        I: AsRef<Path>,
        O: AsRef<Path>,
    {
        let input = input.as_ref();
        let output = output.as_ref();

        let ext = output
            .extension()
            .context("Failed to get extention")?
            .to_str()
            .context("Failed to convert extention to str")?;
        assert_eq!(ext, "mp3");

        let mut command = Command::new("ffmpeg");
        command.arg("-y").arg("-i").arg(input);

        command
            .arg("-metadata")
            .arg(format!("title={}", metadata.title));

        command
            .arg("-metadata")
            .arg(format!("album={}", metadata.album));

        command
            .arg("-metadata")
            .arg(format!("artist={}", self.artists[metadata.artist].name));

        command
            .arg("-metadata")
            .arg(format!("track={}", metadata.track));

        command.arg("-b:a").arg("192k");

        command.arg(output);

        let status = command
            .status()
            .context("Failed to run ffmpeg 'Is ffmpeg installed?'")?;

        Ok(status)
    }
}
