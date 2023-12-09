use anyhow::Context;
use axum::{routing::get, Json, Router};
use serde::Serialize;
use sqlx::{Connection, SqliteConnection};
use tower::ServiceBuilder;
use tower_http::{cors::CorsLayer, services::ServeDir, trace::TraceLayer};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    dotenvy::dotenv().context("Failed to load enviroment variables")?;

    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::DEBUG)
        .init();

    let database_url =
        std::env::var("DATABASE_URL").context("No 'DATABASE_URL' set")?;
    let conn = SqliteConnection::connect(&database_url).await?;

    // Collection Folder Structure:
    //   - public: Served by the server
    //     - tracks: Track Files
    //       - { track_id }.full.flac
    //       - { track_id }.mobile.mp3
    //     - images: Images
    //       - {id}.png
    //
    //   - private: Metadata
    //     - collection.json: Artist[]

    // Artist:
    //  - Id: String
    //  - Name: String
    //  - Picture: String
    //  - Changed: Timestamp
    //  - Albums: Album[]

    // Album:
    //  - Id: String
    //  - Name: String
    //  - CoverArt: String
    //  - Changed: Timestamp
    //  - Tracks: Track[]

    // Track:
    //  - Id: String
    //  - Num: Int
    //  - Name: String
    //  - ArtistID: String
    //  - CoverArt: String
    //  - FileQuality: String
    //  - FileMobile: String
    //  - Changed: Timestamp

    // Routes:
    //   - /api/playlists
    //   - /api/playlists

    let api_router = Router::new().route("/playlists", get(get_all_playlists));

    let app = Router::new()
        .nest("/api", api_router)
        .nest_service("/songs", ServeDir::new("songs"))
        .layer(
            ServiceBuilder::new()
                .layer(TraceLayer::new_for_http())
                .layer(CorsLayer::new()),
        );

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000")
        .await
        .context("Failed to create listener")?;
    axum::serve(listener, app)
        .await
        .context("Failed to serve")?;

    Ok(())
}

#[derive(Serialize, Debug)]
struct Song {
    name: String,
    url: String,
}

#[derive(Serialize, Debug)]
struct Playlist {
    name: String,
    songs: Vec<Song>,
}

async fn get_all_playlists() -> Json<Vec<Playlist>> {
    Json(vec![Playlist {
        name: "Playlist #1".to_string(),
        songs: vec![
            Song {
                name: "Song #1".to_string(),
                url: "http://localhost:3000/songs/song.m4a".to_string(),
            },
            Song {
                name: "Song #2".to_string(),
                url: "http://localhost:3000/songs/".to_string(),
            },
        ],
    }])
}
