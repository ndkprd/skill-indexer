(function () {
  "use strict";

  var prefersReducedMotion = window.matchMedia(
    "(prefers-reduced-motion: reduce)"
  ).matches;

  function fadeIn(el) {
    if (prefersReducedMotion || typeof el.animate !== "function") return;
    el.animate([{ opacity: 0 }, { opacity: 1 }], {
      duration: 150,
      easing: "cubic-bezier(0.25, 1, 0.5, 1)",
    });
  }

  function initInstallCommand() {
    var pre = document.getElementById("install-cmd");
    if (!pre) return;

    var zipPath = pre.getAttribute("data-zip-path");
    if (zipPath) {
      pre.textContent =
        "npx skill-repo-store-install " +
        window.location.origin +
        zipPath +
        " --dest ./skills";
    }

    var copyButton = document.getElementById("copy-cmd");
    if (!copyButton || !navigator.clipboard) return;

    copyButton.addEventListener("click", function () {
      navigator.clipboard.writeText(pre.textContent).then(function () {
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

  function initSearch() {
    var input = document.getElementById("search-input");
    var grid = document.getElementById("card-grid");
    if (!input || !grid) return;

    var noResults = document.getElementById("no-results");
    var cards = Array.prototype.slice.call(grid.querySelectorAll(".card"));

    fetch("/search-index.json")
      .then(function (res) {
        if (!res.ok) throw new Error("search index request failed");
        return res.json();
      })
      .then(function (entries) {
        var fuse = new Fuse(entries, {
          keys: ["name", "description"],
          threshold: 0.35,
        });

        input.addEventListener("input", function () {
          applyFilter(fuse, cards, noResults, input.value.trim());
        });
      })
      .catch(function () {
        input.disabled = true;
        input.placeholder = "Search unavailable";
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

  document.addEventListener("DOMContentLoaded", function () {
    initInstallCommand();
    initSearch();
  });
})();
