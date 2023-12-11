use std::sync::Arc;

use anyhow::Context;
use axum::extract::{
    FromRequestParts, Path, Query, State,
};
use axum::routing::get;
use axum::{Json, Router};
use chrono::{DateTime, Utc};
use error::AppError;
use serde::{Deserialize, Serialize};
use sqlx::{FromRow, PgPool};
use tower::ServiceBuilder;
use tower_http::cors::CorsLayer;
use tower_http::services::ServeDir;
use tower_http::trace::TraceLayer;

use error::Result;

mod error;

const DEFAULT_ARTIST_IMAGE: &str = "/images/artist_default.png";

#[derive(FromRow, Debug)]
struct FetchArtistRes {
    id: String,
    name: String,
    picture: Option<String>,
}

struct DataAccess {
    conn: PgPool,
}

impl DataAccess {
    async fn fetch_all_artists(&self) -> Result<Vec<FetchArtistRes>> {
        sqlx::query_as::<_, FetchArtistRes>("SELECT id, name, picture FROM artists")
            .fetch_all(&self.conn)
            .await
            .map_err(AppError::SqlxError)
    }
}

async fn sync(conn: &PgPool) {
    let collection = dwebble::Collection::get("./test");
    println!("{:#?}", collection.artists());

    // let artists = sqlx::query_as::<_, DbArtist>("SELECT * FROM artists")
    //     .fetch_all(conn)
    //     .await
    //     .unwrap();

    // if let Some(db_artist) = artists.iter().find(|a| a.id == artist.id) {
    //             sqlx::query(
    //                 "
    //                 UPDATE artists SET
    //                     name = $3,
    //                     updated_at = extract(epoch from now())
    //                 WHERE id = $1 AND updated_at < $2",
    //             )
    //             .bind(&db_artist.id)
    //             .bind(artist.changed.0)
    //             .bind(&artist.name)
    //             .execute(conn)
    //             .await
    //             .unwrap();
    //         }

    for artist in collection.artists() {
        sqlx::query(
            "INSERT INTO artists(id, name, picture) VALUES ($1, $2, $3)",
        )
        .bind(&artist.id)
        .bind(&artist.name)
        .bind(artist.picture.as_ref())
        .execute(conn)
        .await
        .unwrap();
    }

    for artist in collection.artists() {
        for album in artist.albums.iter() {
            sqlx::query(
                "INSERT INTO albums(id, name, artist_id) VALUES ($1, $2, $3)",
            )
            .bind(&album.id)
            .bind(&album.name)
            .bind(&artist.id)
            .execute(conn)
            .await
            .unwrap();
        }
    }

    for artist in collection.artists() {
        for album in artist.albums.iter() {
            for track in album.tracks.iter() {
                sqlx::query(
                    "INSERT INTO tracks(
                        id, num, name, artist_id, album_id, file_quality, file_mobile
                    ) VALUES ($1, $2, $3, $4, $5, $6, $7)",
                )
                .bind(&track.id) // $1
                .bind(i32::try_from(track.num).unwrap()) // $2
                .bind(&track.name) // $3
                .bind(&track.artist_id) // $4
                .bind(&album.id) // $5
                .bind(&track.files.quality) // $6
                .bind(&track.files.mobile) // $7
                .execute(conn)
                .await
                .unwrap();
            }
        }
    }
}

struct AppState {
    data_access: DataAccess,
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    dotenvy::dotenv().context("Failed to load enviroment variables")?;

    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::INFO)
        .init();

    let database_url =
        std::env::var("DATABASE_URL").context("No 'DATABASE_URL' set")?;
    let conn = PgPool::connect(&database_url)
        .await
        .context("Failed to connect to postgres database")?;

    sync(&conn).await;

    let state = AppState {
        data_access: DataAccess { conn },
    };

    // Routes:
    //   - /api/search
    //
    //   - /api/artists
    //   - /api/artists/:id
    //   - /api/albums
    //   - /api/albums/:id
    //   - /api/tracks
    //   - /api/tracks/:id
    //   - /api/playlists
    //   - /api/playlists/:id

    async fn stub() {
        unimplemented!();
    }

    let api_router = Router::new()
        .route("/artists", get(get_all_artists))
        .route("/artists/:id", get(get_artist))
        .route("/albums", get(stub))
        .route("/albums/:id", get(stub))
        .route("/tracks", get(stub))
        .route("/tracks/:id", get(stub))
        .route("/playlists", get(stub))
        .route("/playlists/:id", get(stub));

    let state = Arc::new(state);
    let app = Router::new()
        .nest("/api", api_router)
        .nest_service("/songs", ServeDir::new("songs"))
        .layer(
            ServiceBuilder::new()
                .layer(TraceLayer::new_for_http())
                .layer(CorsLayer::new()),
        )
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000")
        .await
        .context("Failed to create listener")?;
    axum::serve(listener, app)
        .await
        .context("Failed to serve")?;

    Ok(())
}

#[derive(Serialize, FromRow)]
struct ArtistRes {
    id: String,
    name: String,
    picture: String,
}

impl From<FetchArtistRes> for ArtistRes {
    fn from(value: FetchArtistRes) -> Self {
        ArtistRes {
            id: value.id,
            name: value.name,
            picture: value.picture.unwrap_or(DEFAULT_ARTIST_IMAGE.to_string()),
        }
    }
}

async fn get_all_artists(
    State(state): State<Arc<AppState>>,
) -> Result<Json<Vec<ArtistRes>>> {
    let artists = state.data_access.fetch_all_artists().await?;

    Ok(Json(artists.into_iter().map(ArtistRes::from).collect()))
}

#[derive(Serialize, FromRow)]
struct DbAlbum {
    id: String,
    name: String,
    created_at: DateTime<Utc>,
}

#[derive(Serialize, FromRow)]
struct DbArtistFull {
    id: String,
    name: String,
}

#[derive(Serialize)]
struct ApiArtist {
    #[serde(flatten)]
    artist: DbArtistFull,
    albums: Vec<DbAlbum>,
}

async fn get_artist(
    State(state): State<Arc<AppState>>,
    Path(id): Path<String>,
) -> Result<Json<ApiArtist>> {
    let artist = sqlx::query_as::<_, DbArtistFull>(
        "SELECT id, name FROM artists WHERE id = $1",
    )
    .bind(&id)
    .fetch_optional(&state.data_access.conn)
    .await
    .map_err(AppError::SqlxError)?;

    let artist = artist.ok_or(AppError::NoArtistWithId(id))?;

    let albums = sqlx::query_as::<_, DbAlbum>(
        "SELECT id, name, created_at FROM albums WHERE artist_id = $1",
    )
    .bind(&artist.id)
    .fetch_all(&state.data_access.conn)
    .await
    .map_err(AppError::SqlxError)?;

    Ok(Json(ApiArtist { artist, albums }))
}
