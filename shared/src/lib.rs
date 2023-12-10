// Collection Folder Structure:
//   - public: Served by the server
//     - tracks: Track Files
//       - { track_id }.full.flac
//       - { track_id }.mobile.mp3
//     - images: Images
//       - {id}.png
//
//    - collection.json

use std::path::{Path, PathBuf};

use metadata::{Album, Artist};

pub mod metadata;

pub struct Collection {
    base: PathBuf,

    metadata: metadata::Collection,
}

impl Collection {
    // TODO(patrik): Collection::get_or_create should return an error
    pub fn get_or_create<P>(path: P) -> Self
    where
        P: AsRef<Path>,
    {
        let path = path.as_ref();

        let create_dir = |base: &Path, sub: &str| {
            let mut path = base.to_path_buf();
            path.push(sub);
            std::fs::create_dir_all(path).unwrap();
        };

        std::fs::create_dir_all(path).unwrap();
        create_dir(path, "tracks");
        create_dir(path, "images");

        let mut collection_json_path = path.to_path_buf();
        collection_json_path.push("collection.json");

        let metadata =
            if let Ok(s) = std::fs::read_to_string(collection_json_path) {
                serde_json::from_str::<metadata::Collection>(&s).unwrap()
            } else {
                metadata::Collection::new(Vec::new())
            };

        Self {
            base: path.to_path_buf(),
            metadata,
        }
    }

    pub fn verify(&self) {}

    pub fn get_or_insert_artist<F>(
        &mut self,
        artist_name: &str,
        insert_fn: F,
    ) -> usize
    where
        F: FnOnce() -> Artist,
    {
        if let Some((index, _)) = self
            .metadata
            .artists
            .iter()
            .enumerate()
            .find(|(_, artist)| artist.name == artist_name)
        {
            index
        } else {
            let artist = insert_fn();
            let index = self.metadata.artists.len();
            self.metadata.artists.push(artist);
            index
        }
    }

    pub fn artist_by_index(&self, artist_index: usize) -> Option<&Artist> {
        self.metadata.artists.get(artist_index)
    }

    pub fn album_with_name_exists(
        &self,
        artist_index: usize,
        album_name: &str,
    ) -> bool {
        if let Some(artist) = self.artist_by_index(artist_index) {
            return artist
                .albums
                .iter()
                .find(|album| album.name == album_name)
                .is_some();
        }

        false
    }

    pub fn add_album(
        &mut self,
        artist_index: usize,
        album: Album,
    ) -> Option<()> {
        let artist = self.metadata.artists.get_mut(artist_index)?;
        // TODO(patrik): Check if album already exists
        artist.albums.push(album);
        Some(())
    }

    pub fn track_dir(&self) -> PathBuf {
        self.base.join("tracks")
    }

    // pub fn copy_track<P>(file: &P) -> Option<()> {
    //     None
    // }

    pub fn save_to_disk(&self) {
        let mut collection_json_path = self.base.clone();
        collection_json_path.push("collection.json");

        let s = serde_json::to_string_pretty(&self.metadata).unwrap();
        std::fs::write(collection_json_path, s).unwrap()
    }
}

pub fn create_id() -> String {
    let cc = cuid2::CuidConstructor::new().with_length(18);
    cc.create_id()
}
