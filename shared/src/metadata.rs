use serde::{Deserialize, Serialize};

#[derive(
    Serialize, Deserialize, Copy, Clone, PartialEq, Eq, PartialOrd, Ord, Debug,
)]
#[repr(transparent)]
pub struct Timestamp(pub i64);

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Collection {
    pub artists: Vec<Artist>,
}

impl Collection {
    pub fn new(artists: Vec<Artist>) -> Self {
        Self { artists }
    }
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Artist {
    pub id: String,
    pub name: String,
    pub picture: Option<String>,
    pub changed: Timestamp,
    pub albums: Vec<Album>,
}

impl Artist {
    pub fn new(id: String, name: String) -> Self {
        Self::new_with_albums(id, name, Vec::new())
    }

    pub fn new_with_albums(
        id: String,
        name: String,
        albums: Vec<Album>,
    ) -> Self {
        Self {
            id,
            name,
            picture: None,
            changed: Timestamp(0),
            albums,
        }
    }

    pub fn set_picture(&mut self, picture: String) {
        self.picture = Some(picture);
    }
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Album {
    pub id: String,
    pub name: String,
    pub cover_art: Option<String>,
    pub changed: Timestamp,
    pub tracks: Vec<Track>,
}

impl Album {
    pub fn new(
        id: String,
        name: String,
    ) -> Self {
        Self::new_with_tracks(id, name, Vec::new())
    }

    pub fn new_with_tracks(
        id: String,
        name: String,
        tracks: Vec<Track>,
    ) -> Self {
        Self {
            id,
            name,
            cover_art: None,
            changed: Timestamp(0),
            tracks,
        }
    }

    pub fn add_track(&mut self, track: Track) -> Option<()> {
        // TODO(patrik): Add checks
        self.tracks.push(track);

        Some(())
    }

    pub fn set_cover_art(&mut self, cover_art: String) {
        self.cover_art = Some(cover_art);
    }
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct TrackFiles {
    pub quality: String,
    pub mobile: String,
}

impl TrackFiles {
    pub fn new(quality: String, mobile: String) -> Self {
        Self { quality, mobile }
    }
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Track {
    pub id: String,
    pub num: usize,
    pub name: String,
    #[serde(rename = "artistId")]
    pub artist_id: String,
    #[serde(rename = "coverArt")]
    pub cover_art: Option<String>,
    pub changed: Timestamp,

    pub files: TrackFiles,
}

impl Track {
    pub fn new(
        id: String,
        num: usize,
        name: String,
        artist_id: String,
        files: TrackFiles,
    ) -> Self {
        Self {
            id,
            num,
            name,
            artist_id,
            cover_art: None,
            changed: Timestamp(0),
            files,
        }
    }

    pub fn set_cover_art(&mut self, cover_art: String) {
        self.cover_art = Some(cover_art);
    }
}
