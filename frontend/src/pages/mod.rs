//! Page components for each route

pub mod auth;
pub mod books;
pub mod dashboard;
pub mod learning;
pub mod review;
pub mod teacher;

pub use auth::{LoginPage, RegisterPage, OAuthCallbackPage};
pub use books::BooksPage;
pub use dashboard::DashboardPage;
pub use learning::LearningSessionPage;
pub use review::ReviewPage;
pub use teacher::TeacherModePage;
