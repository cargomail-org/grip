const registerForm = document.getElementById('registerForm');
const loginForm = document.getElementById('loginForm');

// index.js, auth.js
async function parseJSON(response) {
  if (
    response.status === 204 ||
    response.status === 205 ||
    parseInt(response.headers.get('content-length')) === 0
  ) {
    return null;
  }
  return await response.json();
}

if (registerForm) {
  registerForm.onsubmit = async (e) => {
    e?.preventDefault();

    const form = e.currentTarget;
  
    const formData = {
      username: form.querySelector('input[name="username"]').value,
      password: form.querySelector('input[name="password"]').value,
    };
  
    const confirmation = form.querySelector('input[name="confirmation"]').value;
  
    const alert = form.querySelector('div[name="alert"]')
    if (alert) alert.remove();
  
    if (formData.password != confirmation) {
      form.insertAdjacentHTML(
        "beforeend",
        `<div class="alert alert-warning alert-dismissible fade show" role="alert" name="alert">
          Passwords do NOT match
           <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
        </div>`
      );
      return
    }

    let response;
  
    try {
      const result = await fetch("/api/v1/auth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });
  
      response = await parseJSON(result);
  
      if (result.status != 201) {
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

    window.location.href="/login";  
  };
}

if (loginForm) {
  loginForm.onsubmit = async (e) => {
    e?.preventDefault();

    const form = e.currentTarget;
  
    const formData = {
      username: form.querySelector('input[name="username"]').value,
      password: form.querySelector('input[name="password"]').value,
      rememberMe: form.querySelector('input[name="rememberMe"]').checked,
    };
  
    const alert = form.querySelector('div[name="alert"]');
    if (alert) alert.remove();
    let response;
  
    try {
      const result = await fetch("/api/v1/auth/authenticate", {
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
  
    window.location.href="/";
  };
}


