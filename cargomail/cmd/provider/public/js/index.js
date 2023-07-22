document.addEventListener("DOMContentLoaded", function () {
  composeContent(); // in prod should be cargoesContent
});

function composeContent(e) {
  e?.preventDefault();

  document.getElementById("composeContainer").hidden = false;
  document.getElementById("composeLink").classList.add("active");

  document.getElementById("cargoesContainer").hidden = true;
  document.getElementById("cargoesLink").classList.remove("active");

  document.getElementById("filesContainer").hidden = true;
  document.getElementById("filesLink").classList.remove("active");

  document.getElementById("contactsContainer").hidden = true;
  document.getElementById("contactsLink").classList.remove("active");

  document.getElementById("profileContainer").hidden = true;
  document.getElementById("profileLink").classList.remove("active");

  document.getElementById("composePanel").hidden = false;
  document.getElementById("cargoesPanel").hidden = true;
  document.getElementById("filesPanel").hidden = true;
  document.getElementById("contactsPanel").hidden = true;
  document.getElementById("profilePanel").hidden = true;
}

function cargoesContent(e) {
  e?.preventDefault();

  document.getElementById("composeContainer").hidden = true;
  document.getElementById("composeLink").classList.remove("active");

  document.getElementById("cargoesContainer").hidden = false;
  document.getElementById("cargoesLink").classList.add("active");

  document.getElementById("filesContainer").hidden = true;
  document.getElementById("filesLink").classList.remove("active");

  document.getElementById("contactsContainer").hidden = true;
  document.getElementById("contactsLink").classList.remove("active");

  document.getElementById("profileContainer").hidden = true;
  document.getElementById("profileLink").classList.remove("active");

  document.getElementById("composePanel").hidden = true;
  document.getElementById("cargoesPanel").hidden = false;
  document.getElementById("filesPanel").hidden = true;
  document.getElementById("contactsPanel").hidden = true;
  document.getElementById("profilePanel").hidden = true;
}

function filesContent(e) {
  e?.preventDefault();

  document.getElementById("composeContainer").hidden = true;
  document.getElementById("composeLink").classList.remove("active");

  document.getElementById("cargoesContainer").hidden = true;
  document.getElementById("cargoesLink").classList.remove("active");

  document.getElementById("filesContainer").hidden = false;
  document.getElementById("filesLink").classList.add("active");

  document.getElementById("contactsContainer").hidden = true;
  document.getElementById("contactsLink").classList.remove("active");

  document.getElementById("profileContainer").hidden = true;
  document.getElementById("profileLink").classList.remove("active");

  document.getElementById("composePanel").hidden = true;
  document.getElementById("cargoesPanel").hidden = true;
  document.getElementById("filesPanel").hidden = false;
  document.getElementById("contactsPanel").hidden = true;
  document.getElementById("profilePanel").hidden = true;
}

function contactsContent(e) {
  e?.preventDefault();

  document.getElementById("composeContainer").hidden = true;
  document.getElementById("composeLink").classList.remove("active");

  document.getElementById("cargoesContainer").hidden = true;
  document.getElementById("cargoesLink").classList.remove("active");

  document.getElementById("filesContainer").hidden = true;
  document.getElementById("filesLink").classList.remove("active");

  document.getElementById("contactsContainer").hidden = false;
  document.getElementById("contactsLink").classList.add("active");

  document.getElementById("profileContainer").hidden = true;
  document.getElementById("profileLink").classList.remove("active");

  document.getElementById("composePanel").hidden = true;
  document.getElementById("cargoesPanel").hidden = true;
  document.getElementById("filesPanel").hidden = true;
  document.getElementById("contactsPanel").hidden = false;
  document.getElementById("profilePanel").hidden = true;
}

function profileContent(e) {
  e?.preventDefault();

  document.getElementById("composeContainer").hidden = true;
  document.getElementById("composeLink").classList.remove("active");

  document.getElementById("cargoesContainer").hidden = true;
  document.getElementById("cargoesLink").classList.remove("active");

  document.getElementById("filesContainer").hidden = true;
  document.getElementById("filesLink").classList.remove("active");

  document.getElementById("contactsContainer").hidden = true;
  document.getElementById("contactsLink").classList.remove("active");

  document.getElementById("profileContainer").hidden = false;
  document.getElementById("profileLink").classList.add("active");

  document.getElementById("composePanel").hidden = true;
  document.getElementById("cargoesPanel").hidden = true;
  document.getElementById("filesPanel").hidden = true;
  document.getElementById("contactsPanel").hidden = true;
  document.getElementById("profilePanel").hidden = false;
}

const formatBytes = (bytes, decimals = 2) => {
  if (!+bytes) return "0 Bytes";

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["B", "KB", "MB", "GB", "TB"];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
};

// index.js, auth.js
async function parseJSON(response) {
  if (
    response.status === 204 ||
    response.status === 205 ||
    parseInt(response.headers.get("content-length")) === 0
  ) {
    return null;
  }
  return await response.json();
}

const downloadURI = (formId, uri, name) => {
  (async () => {
    const form = document.getElementById(formId);

    const alert = form.querySelector('div[name="alert"]');
    if (alert) alert.remove();
    let response;

    try {
      const result = await fetch(uri, {
        method: "HEAD",
        headers: {
          Accept: "application/json",
        },
      });

      // response = await parseJSON(result); no response data expected in HEAD method request

      if (result.status != 200) {
        const error = new Error(result.statusText);
        error.response = response;
        throw error;
      }
    } catch (error) {
      let errMessage = "unknown error";
      if (
        error != null &&
        "response" in error &&
        error.response != null &&
        error.response.Err
      ) {
        errMessage =
          error.response.Err.charAt(0).toUpperCase() +
          error.response.Err.slice(1);
      } else if (error != null) {
        errMessage = error.message;
      }

      form.insertAdjacentHTML(
        "beforeend",
        `<div class="alert alert-warning alert-dismissible fade show" role="alert" name="alert">
              ${errMessage}
               <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
         </div>`
      );
      return;
    }

    const link = document.createElement("a");
    link.download = name;
    link.href = uri;
    link.click();
  })();
};
