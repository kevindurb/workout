if (!document.cookie.includes('tz=')) {
  document.cookie = `tz=${Intl.DateTimeFormat().resolvedOptions().timeZone}; path=/; max-age=31536000`;
}
