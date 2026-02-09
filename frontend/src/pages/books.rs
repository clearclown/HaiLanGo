//! Books list and upload page

use leptos::*;

use crate::api::{ApiClient, BookItem};
use crate::components::{Alert, Badge, Button, Card, EmptyState, Input, Spinner};
use crate::layouts::MainLayout;

/// Books list with upload functionality
#[component]
pub fn BooksPage() -> impl IntoView {
    let books = create_rw_signal::<Vec<BookItem>>(vec![]);
    let loading = create_rw_signal(true);
    let show_upload = create_rw_signal(false);
    let upload_error = create_rw_signal::<Option<String>>(None);
    let upload_loading = create_rw_signal(false);

    // Upload form fields
    let title = create_rw_signal(String::new());
    let author = create_rw_signal(String::new());
    let language = create_rw_signal("zh".to_string());

    // Fetch books
    create_effect(move |_| {
        spawn_local(async move {
            match ApiClient::get_books().await {
                Ok(book_list) => books.set(book_list),
                Err(e) => log::error!("Failed to fetch books: {}", e),
            }
            loading.set(false);
        });
    });

    let on_upload = move |ev: web_sys::SubmitEvent| {
        ev.prevent_default();
        upload_loading.set(true);
        upload_error.set(None);

        let t = title.get();
        let a = author.get();
        let l = language.get();

        spawn_local(async move {
            match ApiClient::upload_book(&t, &a, &l).await {
                Ok(book) => {
                    books.update(|b| b.push(book));
                    show_upload.set(false);
                    title.set(String::new());
                    author.set(String::new());
                }
                Err(e) => upload_error.set(Some(e)),
            }
            upload_loading.set(false);
        });
    };

    view! {
        <MainLayout>
            <div class="books-page">
                <div class="page-header">
                    <h1>"My Books"</h1>
                    <Button
                        text="Add Book".to_string()
                        variant="primary"
                        on_click=Callback::new(move |_| show_upload.set(true))
                    />
                </div>

                // Upload form (shown as inline card)
                {move || if show_upload.get() {
                    view! {
                        <Card title="Add New Book".to_string() class="upload-card">
                            <form class="upload-form" on:submit=on_upload>
                                {move || upload_error.get().map(|msg| view! {
                                    <Alert message=msg variant="danger" />
                                })}

                                <Input
                                    label="Title"
                                    name="book_title"
                                    placeholder="Book title".to_string()
                                    value=title
                                    required=true
                                />
                                <Input
                                    label="Author"
                                    name="book_author"
                                    placeholder="Author name".to_string()
                                    value=author
                                    required=true
                                />

                                <div class="form-group">
                                    <label class="form-label" for="book_language">"Language"</label>
                                    <select
                                        id="book_language"
                                        class="form-input"
                                        on:change=move |ev| language.set(event_target_value(&ev))
                                    >
                                        <option value="zh" selected=true>"Chinese"</option>
                                        <option value="en">"English"</option>
                                        <option value="ja">"Japanese"</option>
                                        <option value="ko">"Korean"</option>
                                        <option value="es">"Spanish"</option>
                                        <option value="fr">"French"</option>
                                        <option value="de">"German"</option>
                                    </select>
                                </div>

                                <div class="form-actions">
                                    <Button
                                        text=Signal::derive(move || {
                                            if upload_loading.get() { "Adding...".to_string() }
                                            else { "Add Book".to_string() }
                                        })
                                        variant="primary"
                                        disabled=upload_loading
                                    />
                                    <Button
                                        text="Cancel".to_string()
                                        variant="secondary"
                                        on_click=Callback::new(move |_| show_upload.set(false))
                                    />
                                </div>
                            </form>
                        </Card>
                    }.into_view()
                } else {
                    view! {}.into_view()
                }}

                // Books list
                {move || if loading.get() {
                    view! { <Spinner /> }.into_view()
                } else if books.get().is_empty() {
                    view! {
                        <EmptyState
                            message="No books yet. Add your first book to start learning!"
                            icon="\u{1F4DA}".to_string()
                        />
                    }.into_view()
                } else {
                    view! {
                        <div class="books-grid">
                            <For
                                each=move || books.get()
                                key=|book| book.id.clone()
                                children=move |book| {
                                    let learn_href = format!("/learn/{}", book.id);
                                    let teacher_href = format!("/teacher?book={}", book.id);
                                    view! {
                                        <div class="book-card">
                                            <Card>
                                                <div class="book-card-content">
                                                    <h3 class="book-title">{book.title.clone()}</h3>
                                                    <p class="book-author">{book.author.clone()}</p>
                                                    <div class="book-info">
                                                        <Badge text=book.language.clone() variant="info" />
                                                        <span class="book-pages">
                                                            {format!("{} pages", book.total_pages)}
                                                        </span>
                                                    </div>
                                                    <div class="book-progress-bar">
                                                        <div
                                                            class="book-progress-fill"
                                                            style=format!("width: {}%", book.progress * 100.0)
                                                        />
                                                    </div>
                                                    <div class="book-actions">
                                                        <a href=learn_href>
                                                            <Button text="Learn".to_string() variant="primary" />
                                                        </a>
                                                        <a href=teacher_href>
                                                            <Button text="Teacher".to_string() variant="secondary" />
                                                        </a>
                                                    </div>
                                                </div>
                                            </Card>
                                        </div>
                                    }
                                }
                            />
                        </div>
                    }.into_view()
                }}
            </div>
        </MainLayout>
    }
}
