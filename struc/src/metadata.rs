use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Debug)]
pub struct TrackMetadata {
    #[serde(rename = "trackNum")]
    pub track_num: usize,
    pub name: String,
    pub artist_id: String,

    pub quality_version: String,
    pub mobile_version: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct AlbumMetadata {
    pub id: String,
    pub name: String,
    #[serde(rename = "artistId")]
    pub artist_id: Option<String>,
    #[serde(rename = "coverArt")]
    pub cover_art: String,

    pub tracks: Vec<TrackMetadata>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct ArtistMetadata {
    pub id: String,
    pub name: String,
    pub picture: String,
}
