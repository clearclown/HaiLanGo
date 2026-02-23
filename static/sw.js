/**
 * HaiLanGo Service Worker
 * PWA offline support: caches static assets and TTS audio files
 */

const CACHE_VERSION = 'v1';
const STATIC_CACHE = `hailango-static-${CACHE_VERSION}`;
const AUDIO_CACHE = `hailango-audio-${CACHE_VERSION}`;

const STATIC_ASSETS = [
    '/',
    '/styles/main.css',
    '/favicon.png',
    '/manifest.json',
];

// Install: pre-cache static assets
self.addEventListener('install', (event) => {
    event.waitUntil(
        caches.open(STATIC_CACHE).then((cache) => cache.addAll(STATIC_ASSETS))
    );
    self.skipWaiting();
});

// Activate: clean up old caches
self.addEventListener('activate', (event) => {
    event.waitUntil(
        caches.keys().then((keys) =>
            Promise.all(
                keys
                    .filter((k) => k !== STATIC_CACHE && k !== AUDIO_CACHE)
                    .map((k) => caches.delete(k))
            )
        )
    );
    self.clients.claim();
});

// Fetch: serve from cache with network fallback
self.addEventListener('fetch', (event) => {
    const { request } = event;
    const url = new URL(request.url);

    // Cache TTS audio responses (offline teacher mode support)
    if (url.pathname.startsWith('/api/tts/')) {
        event.respondWith(
            caches.open(AUDIO_CACHE).then(async (cache) => {
                const cached = await cache.match(request);
                if (cached) return cached;
                const response = await fetch(request);
                if (response.ok) cache.put(request, response.clone());
                return response;
            })
        );
        return;
    }

    // For API requests: network only (no caching of dynamic data)
    if (url.pathname.startsWith('/api/')) {
        event.respondWith(fetch(request));
        return;
    }

    // For static assets: cache first
    event.respondWith(
        caches.match(request).then((cached) => cached || fetch(request))
    );
});
