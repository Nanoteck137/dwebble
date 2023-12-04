use axum::{routing::get, Json, Router};
use serde::Serialize;
use tower::ServiceBuilder;
use tower_http::{cors::CorsLayer, services::ServeDir, trace::TraceLayer};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::DEBUG)
        .init();

    let api_router = Router::new().route("/playlists", get(get_all_playlists));

    let app = Router::new()
        .nest("/api", api_router)
        .nest_service("/songs", ServeDir::new("songs"))
        .layer(
            ServiceBuilder::new()
                .layer(TraceLayer::new_for_http())
                .layer(CorsLayer::new()),
        );

    let listener =
        tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();

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
