use axum::{http::StatusCode, response::IntoResponse, Json};
use serde_json::json;
use thiserror::Error;
use tracing::Level;

pub type Result<T> = core::result::Result<T, AppError>;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("No artist with id: '{0}'")]
    NoArtistWithId(String),

    #[error("Internal server error")]
    SqlxError(sqlx::Error),
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
        };

        let body = Json(json!({
            "error": error_message,
        }));

        (status, body).into_response()
    }
}
