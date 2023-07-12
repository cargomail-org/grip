document.addEventListener("DOMContentLoaded", function() {
    collectionsContent();
  });

function composeContent(e) {
    document.getElementById("compose-container").hidden = false;
    document.getElementById("compose-link").classList.add("active");

    document.getElementById("collections-container").hidden = true;
    document.getElementById("collections-link").classList.remove("active");

    document.getElementById("files-container").hidden = true;
    document.getElementById("files-link").classList.remove("active");
}

function collectionsContent(e) {
    document.getElementById("compose-container").hidden = true;
    document.getElementById("compose-link").classList.remove("active");

    document.getElementById("collections-container").hidden = false;
    document.getElementById("collections-link").classList.add("active");

    document.getElementById("files-container").hidden = true;
    document.getElementById("files-link").classList.remove("active");
}

function filesContent(e) {
    document.getElementById("compose-container").hidden = true;
    document.getElementById("compose-link").classList.remove("active");
    
    document.getElementById("collections-container").hidden = true;
    document.getElementById("collections-link").classList.remove("active");
    
    document.getElementById("files-container").hidden = false;
    document.getElementById("files-link").classList.add("active");
}