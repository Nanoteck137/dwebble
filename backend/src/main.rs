use std::sync::Arc;

use anyhow::Context;
use axum::extract::rejection::{JsonRejection, QueryRejection};
use axum::extract::{
    FromRequest, FromRequestParts, Path, Query, Request, State,
};
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::routing::get;
use axum::{Json, Router};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use serde_json::json;
use sqlx::{FromRow, PgPool};
use thiserror::Error;
use tower::ServiceBuilder;
use tower_http::cors::CorsLayer;
use tower_http::services::ServeDir;
use tower_http::trace::TraceLayer;
use tracing::Level;

type Result<T> = core::result::Result<T, AppError>;

const DEFAULT_ARTIST_IMAGE: &str = "/images/artist_default.png";

#[derive(FromRow, Debug)]
struct FetchArtistRes {
    id: String,
    name: String,
    picture: Option<String>,
    created_at: DateTime<Utc>,
    updated_at: i64,
}

struct DataAccess {
    conn: PgPool,
}

impl DataAccess {
    async fn fetch_all_artists(&self) -> Result<Vec<FetchArtistRes>> {
        sqlx::query_as::<_, FetchArtistRes>("SELECT * FROM artists")
            .fetch_all(&self.conn)
            .await
            .map_err(AppError::SqlxError)
    }
}

#[derive(FromRow)]
struct DbArtist {
    id: String,
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

    let api_router = Router::new()
        .route("/artists", get(get_all_artists))
        .route("/artists/:id", get(get_artist))
        .route("/artists/search", get(search_for_artist));

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

// #[derive(Deserialize, Validate, Debug)]
// struct Test {
//     #[validate(length(min = 5))]
//     field: String,
//     field2: usize,
// }
//
// async fn test(ValidatedJson(body): ValidatedJson<Test>) {
//     println!("Body: {:#?}", body);
// }
//
// #[derive(Debug, Clone, Copy, Default)]
// pub struct ValidatedJson<T>(pub T);
//
// #[async_trait]
// impl<T, S> FromRequest<S> for ValidatedJson<T>
// where
//     T: DeserializeOwned + Validate,
//     S: Send + Sync,
//     Json<T>: FromRequest<S, Rejection = JsonRejection>,
// {
//     type Rejection = String;
//
//     async fn from_request(
//         req: Request,
//         state: &S,
//     ) -> Result<Self, Self::Rejection> {
//         let Json(value) = Json::<T>::from_request(req, state).await.unwrap();
//         value.validate().unwrap();
//         Ok(ValidatedJson(value))
//     }
// }

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

#[derive(FromRequestParts)]
#[from_request(via(Query), rejection(AppError))]
pub struct QueryExtractor<T>(pub T);

#[derive(Deserialize, Debug)]
struct SearchQueryParams {
    q: String,
}

async fn search_for_artist(
    State(state): State<Arc<AppState>>,
    QueryExtractor(query_params): QueryExtractor<SearchQueryParams>,
) -> Result<Json<Vec<ArtistRes>>> {
    sqlx::query_as::<_, ArtistRes>(
        "SELECT id, name FROM artists WHERE name LIKE $1",
    )
    .bind(format!("%{}%", query_params.q))
    .fetch_all(&state.data_access.conn)
    .await
    .map(|a| Json(a))
    .map_err(AppError::SqlxError)
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

#[derive(Error, Debug)]
pub enum AppError {
    #[error("No artist with id: '{0}'")]
    NoArtistWithId(String),

    #[error("Internal server error")]
    SqlxError(sqlx::Error),

    #[error("Query error: {0}")]
    QueryError(#[from] QueryRejection),
}

impl IntoResponse for AppError {
    #[tracing::instrument]
    fn into_response(self) -> axum::response::Response {
        let (status, error_message) = match self {
            AppError::NoArtistWithId(id) => {
                (StatusCode::NOT_FOUND, format!("No artist with id '{}'", id))
            }

            AppError::SqlxError(e) => {
                tracing::event!(Level::ERROR, "Sqlx Error: {e}");
                (
                    StatusCode::INTERNAL_SERVER_ERROR,
                    "Internal server error".to_string(),
                )
            }

            AppError::QueryError(e) => {
                tracing::event!(Level::ERROR, "Query Error: {e}");
                (e.status(), e.body_text())
            }
        };

        let body = Json(json!({
            "error": error_message,
        }));

        (status, body).into_response()
    }
}
