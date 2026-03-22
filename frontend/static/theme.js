(function () {
  var storageKey = 'buddy-theme';
  var root = document.documentElement;
  var mediaQuery = typeof window.matchMedia === 'function'
    ? window.matchMedia('(prefers-color-scheme: dark)')
    : null;
  var followsSystemTheme = false;

  function readStoredTheme() {
    try {
      var theme = window.localStorage.getItem(storageKey);
      return theme === 'dark' || theme === 'light' ? theme : null;
    } catch (error) {
      return null;
    }
  }

  function writeStoredTheme(theme) {
    try {
      window.localStorage.setItem(storageKey, theme);
    } catch (error) {
      // Ignore storage failures and still update the page theme.
    }
  }

  function clearStoredTheme() {
    try {
      window.localStorage.removeItem(storageKey);
    } catch (error) {
      // Ignore storage failures and keep using the system preference.
    }
  }

  function getSystemTheme() {
    return mediaQuery && mediaQuery.matches ? 'dark' : 'light';
  }

  function currentTheme() {
    return root.getAttribute('data-theme') || resolveTheme();
  }

  function updateToggleState(toggle) {
    if (!toggle) {
      return;
    }

    var isDark = currentTheme() === 'dark';
    var nextThemeLabel = isDark ? 'Switch to light mode' : 'Switch to dark mode';

    toggle.setAttribute('aria-pressed', isDark ? 'true' : 'false');
    toggle.setAttribute('aria-label', nextThemeLabel);
    toggle.setAttribute('title', nextThemeLabel);
  }

  function refreshToggles() {
    document.querySelectorAll('[data-theme-toggle]').forEach(updateToggleState);
  }

  function applyTheme(theme) {
    root.setAttribute('data-theme', theme);
    refreshToggles();
    return theme;
  }

  function resolveTheme() {
    var storedTheme = readStoredTheme();
    followsSystemTheme = storedTheme === null;
    return storedTheme || getSystemTheme();
  }

  function setTheme(theme) {
    var normalizedTheme = theme === 'dark' ? 'dark' : 'light';
    followsSystemTheme = false;
    writeStoredTheme(normalizedTheme);
    return applyTheme(normalizedTheme);
  }

  function clearThemePreference() {
    followsSystemTheme = true;
    clearStoredTheme();
    return applyTheme(getSystemTheme());
  }

  function toggleTheme() {
    return setTheme(currentTheme() === 'dark' ? 'light' : 'dark');
  }

  function bindToggle(toggle) {
    if (!toggle || toggle.dataset.themeToggleBound === 'true') {
      return;
    }

    toggle.addEventListener('click', function (event) {
      event.preventDefault();
      toggleTheme();
    });

    toggle.dataset.themeToggleBound = 'true';
    updateToggleState(toggle);
  }

  function attachToggle(target) {
    if (!target) {
      return api;
    }

    if (typeof target === 'string') {
      document.querySelectorAll(target).forEach(bindToggle);
      return api;
    }

    if (typeof target.length === 'number' && typeof target !== 'function') {
      Array.prototype.forEach.call(target, bindToggle);
      return api;
    }

    bindToggle(target);
    return api;
  }

  var api = {
    getTheme: currentTheme,
    setTheme: setTheme,
    toggleTheme: toggleTheme,
    clearThemePreference: clearThemePreference,
    attachToggle: attachToggle
  };

  applyTheme(resolveTheme());
  window.BuddyTheme = api;

  if (mediaQuery) {
    var handleSystemThemeChange = function () {
      if (followsSystemTheme) {
        applyTheme(getSystemTheme());
      }
    };

    if (typeof mediaQuery.addEventListener === 'function') {
      mediaQuery.addEventListener('change', handleSystemThemeChange);
    } else if (typeof mediaQuery.addListener === 'function') {
      mediaQuery.addListener(handleSystemThemeChange);
    }
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function () {
      attachToggle('[data-theme-toggle]');
    });
  } else {
    attachToggle('[data-theme-toggle]');
  }
})();
