(function () {
  "use strict";

  var prefersReducedMotion = window.matchMedia(
    "(prefers-reduced-motion: reduce)"
  ).matches;
  var panelCloseDelay = prefersReducedMotion ? 0 : 190;

  var skillsByDir = {};
  var activeTrigger = null;

  // ---- Theme ----

  function getEffectiveTheme() {
    var stored = readStoredTheme();
    if (stored) return stored;
    return window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
  }

  function readStoredTheme() {
    try {
      var t = localStorage.getItem("theme");
      return t === "light" || t === "dark" ? t : null;
    } catch (e) {
      return null;
    }
  }

  function writeStoredTheme(theme) {
    try {
      localStorage.setItem("theme", theme);
    } catch (e) {}
  }

  function paintToggle(theme) {
    var toggle = document.getElementById("theme-toggle");
    if (!toggle) return;
    toggle.classList.toggle("is-dark", theme === "dark");
    toggle.setAttribute("aria-pressed", theme === "dark" ? "true" : "false");
  }

  function setTheme(theme) {
    document.documentElement.setAttribute("data-theme", theme);
    writeStoredTheme(theme);
    paintToggle(theme);
  }

  function initTheme() {
    paintToggle(getEffectiveTheme());

    var toggle = document.getElementById("theme-toggle");
    if (toggle) {
      toggle.addEventListener("click", function () {
        setTheme(getEffectiveTheme() === "dark" ? "light" : "dark");
      });
    }

    window
      .matchMedia("(prefers-color-scheme: dark)")
      .addEventListener("change", function () {
        if (!readStoredTheme()) paintToggle(getEffectiveTheme());
      });
  }

  // ---- Detail panel ----

  function formatMetadataValue(v) {
    if (Array.isArray(v)) return v.join(", ");
    if (v === null || v === undefined) return "";
    return String(v);
  }

  function populatePanel(entry) {
    document.getElementById("panel-name").textContent = entry.name;
    document.getElementById("panel-dirname").textContent = entry.dirName;
    document.getElementById("panel-description").textContent =
      entry.description;

    var section = document.getElementById("panel-metadata-section");
    var body = document.getElementById("panel-metadata-body");
    var keys = Object.keys(entry.metadata || {}).sort();
    body.innerHTML = "";
    keys.forEach(function (key) {
      var tr = document.createElement("tr");
      var th = document.createElement("th");
      th.scope = "row";
      th.textContent = key;
      var td = document.createElement("td");
      td.textContent = formatMetadataValue(entry.metadata[key]);
      tr.appendChild(th);
      tr.appendChild(td);
      body.appendChild(tr);
    });
    section.hidden = keys.length === 0;

    var download = document.getElementById("panel-download");
    download.href = entry.zipPath;

    var pre = document.getElementById("install-cmd");
    pre.textContent =
      "npx skill-repo-store-install " +
      window.location.origin +
      entry.zipPath +
      " --dest ./skills";
  }

  function showPanelUI(entry, trigger) {
    populatePanel(entry);

    var panel = document.getElementById("skill-panel");
    var backdrop = document.getElementById("panel-backdrop");
    panel.hidden = false;
    backdrop.hidden = false;
    document.body.classList.add("panel-open");

    // Force layout so the transform/opacity transition actually runs.
    // eslint-disable-next-line no-unused-expressions
    panel.offsetHeight;
    panel.classList.add("is-open");
    backdrop.classList.add("is-open");
    panel.classList.remove("is-closing");

    setInertOutsidePanel(true);
    activeTrigger = trigger || document.activeElement;
    document.getElementById("panel-close").focus();

    document.addEventListener("keydown", onPanelKeydown, true);
  }

  function closePanelUI() {
    var panel = document.getElementById("skill-panel");
    var backdrop = document.getElementById("panel-backdrop");
    if (panel.hidden) return;

    panel.classList.add("is-closing");
    panel.classList.remove("is-open");
    backdrop.classList.remove("is-open");
    document.body.classList.remove("panel-open");
    setInertOutsidePanel(false);
    document.removeEventListener("keydown", onPanelKeydown, true);

    window.setTimeout(function () {
      panel.hidden = true;
      backdrop.hidden = true;
      panel.classList.remove("is-closing");
    }, panelCloseDelay);

    if (activeTrigger && document.contains(activeTrigger)) {
      activeTrigger.focus();
    }
    activeTrigger = null;
  }

  function setInertOutsidePanel(isInert) {
    ["header", "main", "footer"].forEach(function (selector) {
      var el = document.querySelector("." + siteClassFor(selector));
      if (!el) return;
      if (isInert) el.setAttribute("inert", "");
      else el.removeAttribute("inert");
    });
  }

  function siteClassFor(tag) {
    return (
      { header: "site-header", main: "site-main", footer: "site-footer" }[
        tag
      ] || tag
    );
  }

  function onPanelKeydown(event) {
    if (event.key === "Escape") {
      event.preventDefault();
      requestClose();
      return;
    }
    if (event.key !== "Tab") return;

    var panel = document.getElementById("skill-panel");
    var focusables = panel.querySelectorAll(
      'a[href], button:not([disabled]), [tabindex]:not([tabindex="-1"])'
    );
    if (focusables.length === 0) return;

    var first = focusables[0];
    var last = focusables[focusables.length - 1];

    if (event.shiftKey && document.activeElement === first) {
      event.preventDefault();
      last.focus();
    } else if (!event.shiftKey && document.activeElement === last) {
      event.preventDefault();
      first.focus();
    }
  }

  function openPanel(dirName, trigger) {
    if (!skillsByDir[dirName]) return;
    if (window.location.hash.slice(1) !== dirName) {
      history.pushState({ panel: dirName }, "", "#" + dirName);
    }
    showPanelUI(skillsByDir[dirName], trigger);
  }

  function requestClose() {
    if (history.state && history.state.panel) {
      history.back();
    } else {
      closePanelUI();
    }
  }

  function syncPanelToHash() {
    var dirName = window.location.hash.slice(1);
    if (dirName && skillsByDir[dirName]) {
      showPanelUI(skillsByDir[dirName], null);
    } else {
      closePanelUI();
    }
  }

  function initPanel(entries) {
    entries.forEach(function (entry) {
      skillsByDir[entry.dirName] = entry;
    });

    document
      .getElementById("card-grid")
      .addEventListener("click", function (event) {
        var card = event.target.closest(".card");
        if (!card) return;
        event.preventDefault();
        openPanel(card.getAttribute("data-skill"), card);
      });

    document
      .getElementById("panel-close")
      .addEventListener("click", requestClose);
    document
      .getElementById("panel-backdrop")
      .addEventListener("click", requestClose);

    window.addEventListener("popstate", syncPanelToHash);

    initInstallCopyButton();
    syncPanelToHash();
  }

  function initInstallCopyButton() {
    var copyButton = document.getElementById("copy-cmd");
    if (!copyButton || !navigator.clipboard) return;

    copyButton.addEventListener("click", function () {
      var text = document.getElementById("install-cmd").textContent;
      navigator.clipboard.writeText(text).then(function () {
        var original = copyButton.textContent;
        copyButton.textContent = "Copied";
        copyButton.classList.add("is-copied");
        window.setTimeout(function () {
          copyButton.textContent = original;
          copyButton.classList.remove("is-copied");
        }, 1600);
      });
    });
  }

  // ---- Search ----

  function fadeIn(el) {
    if (prefersReducedMotion || typeof el.animate !== "function") return;
    el.animate([{ opacity: 0 }, { opacity: 1 }], {
      duration: 150,
      easing: "cubic-bezier(0.25, 1, 0.5, 1)",
    });
  }

  function initSearch(entries) {
    var input = document.getElementById("search-input");
    var grid = document.getElementById("card-grid");
    if (!input || !grid) return;

    var noResults = document.getElementById("no-results");
    var cards = Array.prototype.slice.call(grid.querySelectorAll(".card"));
    var fuse = new Fuse(entries, {
      keys: ["name", "description"],
      threshold: 0.35,
    });

    input.addEventListener("input", function () {
      applyFilter(fuse, cards, noResults, input.value.trim());
    });
  }

  function applyFilter(fuse, cards, noResults, query) {
    if (!query) {
      cards.forEach(function (card) {
        showCard(card);
      });
      setNoResults(noResults, "");
      return;
    }

    var matched = new Set(
      fuse.search(query).map(function (result) {
        return result.item.dirName;
      })
    );

    var visibleCount = 0;
    cards.forEach(function (card) {
      var dirName = card.getAttribute("data-skill");
      if (matched.has(dirName)) {
        showCard(card);
        visibleCount += 1;
      } else {
        card.hidden = true;
      }
    });

    setNoResults(
      noResults,
      visibleCount === 0
        ? "No skills match “" + query + "”. Try a different search term."
        : ""
    );
  }

  function showCard(card) {
    var wasHidden = card.hidden;
    card.hidden = false;
    if (wasHidden) fadeIn(card);
  }

  function setNoResults(el, message) {
    if (!el) return;
    el.textContent = message;
    el.hidden = message === "";
  }

  // ---- Boot ----

  document.addEventListener("DOMContentLoaded", function () {
    initTheme();

    fetch("/search-index.json")
      .then(function (res) {
        if (!res.ok) throw new Error("search index request failed");
        return res.json();
      })
      .then(function (entries) {
        initPanel(entries);
        initSearch(entries);
      })
      .catch(function () {
        var input = document.getElementById("search-input");
        if (input) {
          input.disabled = true;
          input.placeholder = "Search unavailable";
        }
        var notice = document.getElementById("load-error");
        if (notice) notice.hidden = false;
      });
  });
})();
