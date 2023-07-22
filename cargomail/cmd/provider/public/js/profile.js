const profileForm = document.getElementById("profileForm");

profileForm.onsubmit = async (e) => {
  e.preventDefault();

  const form = e.currentTarget;

  const formData = {
    firstname: form.querySelector('input[name="firstname"]').value,
    lastname: form.querySelector('input[name="lastname"]').value,
  };

  const alert = form.querySelector('div[name="alert"]');
  if (alert) alert.remove();
  let response;

  try {
    const result = await fetch("/api/v1/user/profile", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
    });

    response = await parseJSON(result);

    if (result.status != 200) {
        const error = new Error(result.statusText);
        error.response = response;
        throw error;
    }
  } catch (error) {
    let errMessage = "unknown error";
    if (error != null && "response" in error && error.response != null && error.response.Err) {
        errMessage = error.response.Err.charAt(0).toUpperCase() + error.response.Err.slice(1);
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
    return
  }
};
