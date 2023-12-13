use crate::error::{AppError, Result};
use serde::Serialize;
use sqlx::{postgres::PgRow, FromRow, PgPool, Row};

#[derive(Serialize, Debug)]
pub struct Artist {
    pub id: String,
    pub name: String,
    pub picture: String,
}

#[derive(Serialize, Debug)]
pub struct Album {
    pub id: String,
    pub name: String,
    pub picture: String,
}

#[derive(Serialize, Debug)]
pub struct Track {
    pub id: String,
    pub num: i32,
    pub name: String,
    pub artist_id: String,
    pub album_id: String,

    pub album_name: String,

    pub file_quality: String,
    pub file_mobile: String,
}

pub struct DataAccess {
    conn: PgPool,
}

impl DataAccess {
    pub fn new(conn: PgPool) -> Self {
        Self { conn }
    }

    fn map_artist(row: PgRow) -> sqlx::Result<Artist> {
        Ok(Artist {
            id: row.try_get("id")?,
            name: row.try_get("name")?,
            picture: row
                .try_get::<Option<String>, _>("picture")?
                .unwrap_or(crate::DEFAULT_ARTIST_IMAGE.to_string()),
        })
    }

    fn map_album(row: PgRow) -> sqlx::Result<Album> {
        Ok(Album {
            id: row.try_get("id")?,
            name: row.try_get("name")?,
            picture: row
                .try_get::<Option<String>, _>("picture")?
                .unwrap_or(crate::DEFAULT_ALBUM_IMAGE.to_string()),
        })
    }

    fn map_track(row: PgRow) -> sqlx::Result<Track> {
        Ok(Track {
            id: row.try_get("id")?,
            num: row.try_get("num")?,
            name: row.try_get("name")?,
            artist_id: row.try_get("artist_id")?,
            album_id: row.try_get("album_id")?,
            album_name: row.try_get("album_name")?,
            file_quality: row.try_get("file_quality")?,
            file_mobile: row.try_get("file_mobile")?,
        })
    }

    pub async fn fetch_all_artists(&self) -> Result<Vec<Artist>> {
        sqlx::query("SELECT id, name, picture FROM artists")
            .try_map(Self::map_artist)
            .fetch_all(&self.conn)
            .await
            .map_err(AppError::SqlxError)
    }

    pub async fn fetch_single_artist(&self, id: &str) -> Result<Artist> {
        sqlx::query(
            "SELECT id, name, picture FROM artists WHERE id = $1",
        )
        .bind(id)
        .try_map(Self::map_artist)
        .fetch_optional(&self.conn)
        .await
        .map_err(AppError::SqlxError)?
        .ok_or_else(|| AppError::NoArtistWithId(id.to_string()))
    }

    pub async fn fetch_albums_by_artist(
        &self,
        id: &str,
    ) -> Result<Vec<Album>> {
        sqlx::query(
            "SELECT id, name, picture FROM albums WHERE artist_id = $1",
        )
        .bind(id)
        .try_map(Self::map_album)
        .fetch_all(&self.conn)
        .await
        .map_err(AppError::SqlxError)
    }

    pub async fn fetch_single_album(&self, id: &str) -> Result<Album> {
        sqlx::query(
            "SELECT id, name, picture FROM albums WHERE id = $1",
        )
        .bind(id)
        .try_map(Self::map_album)
        .fetch_optional(&self.conn)
        .await
        .map_err(AppError::SqlxError)?
        .ok_or_else(|| AppError::NoAlbumWithId(id.to_string()))
    }

    pub async fn fetch_tracks_by_album(&self, id: &str) -> Result<Vec<Track>> {
        sqlx::query(
            "
            SELECT 
                tracks.id as id, 
                tracks.num as num, 
                tracks.name as name, 
                tracks.artist_id as artist_id, 
                tracks.album_id as album_id, 
                tracks.file_quality as file_quality, 
                tracks.file_mobile as file_mobile, 
                albums.name as album_name 
            FROM tracks 
            JOIN albums ON albums.id = tracks.album_id 
            WHERE album_id = $1",
        )
        .bind(id)
        .try_map(Self::map_track)
        .fetch_all(&self.conn)
        .await
        .map_err(AppError::SqlxError)
    }
}
