use anyhow::Context;
use axum::routing::get;
use axum::{Json, Router};
use serde::Serialize;
use sqlx::sqlite::SqliteRow;
use sqlx::types::chrono::Utc;
use sqlx::{Connection, Row, SqliteConnection};
use tower::ServiceBuilder;
use tower_http::cors::CorsLayer;
use tower_http::services::ServeDir;
use tower_http::trace::TraceLayer;

async fn sync(conn: &mut SqliteConnection) {
    let collection = dwebble::Collection::get("./test");
    println!("{:#?}", collection.artists());

    for artist in collection.artists() {
        sqlx::query("INSERT INTO artists(id, name, test) VALUES (?, ?, ?)")
            .bind(&artist.id)
            .bind(&artist.name)
            .bind(1702210874)
            .execute(&mut *conn)
            .await
            .unwrap();

        for album in artist.albums.iter() {
            sqlx::query(
                "INSERT INTO albums(id, name, artist_id) VALUES (?, ?, ?)",
            )
            .bind(&album.id)
            .bind(&album.name)
            .bind(&artist.id)
            .execute(&mut *conn)
            .await
            .unwrap();

            for track in album.tracks.iter() {
                sqlx::query(
                    "INSERT INTO tracks(id, name, album_id) VALUES (?, ?, ?)",
                )
                .bind(&track.id)
                .bind(&track.name)
                .bind(&album.id)
                .execute(&mut *conn)
                .await
                .unwrap();
            }
        }
    }

    #[derive(Debug)]
    struct Test {
        name: String,
        created_at: String,
        test: i64,
    }

    let list = sqlx::query("SELECT * FROM artists")
        .map(|row: SqliteRow| Test {
            name: row.get("name"),
            created_at: row.get("created_at"),
            test: row.get("test"),
        })
        .fetch_all(&mut *conn)
        .await
        .unwrap();
    println!("{:#?}", list);

    // let typ = sqlx::query!("SELECT * FROM artists")
    //     .fetch_one(&mut *conn)
    //     .await
    //     .unwrap();
    // println!("Typ: {:?}", typ.name);

    // let list = sqlx::query!("SELECT * FROM artists").fetch_all(&mut *conn).await.unwrap();
    // println!("List: {:#?}", list);
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    dotenvy::dotenv().context("Failed to load enviroment variables")?;

    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::DEBUG)
        .init();

    let database_url =
        std::env::var("DATABASE_URL").context("No 'DATABASE_URL' set")?;
    let mut conn = SqliteConnection::connect(&database_url).await?;

    // let x = sqlx::query!("SELECT DATE()");

    sync(&mut conn).await;

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
