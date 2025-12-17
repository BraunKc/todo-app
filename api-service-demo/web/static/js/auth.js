document.addEventListener("keydown", function(event) {
  if (event.key == "Tab") {
    event.preventDefault();
  }
});


let authCard = document.getElementById("authCard");
let loginForm = document.getElementById("loginForm");
document.getElementById("toRegistration").addEventListener("click", () => {
    authCard.classList.remove("activate-login");
    registrationForm.removeAttribute("style");
    authCard.classList.add("activate-registration");

    function handle(e) {
        if (e.propertyName != "transform" || e.target != loginForm) return;
        loginForm.style.transition = "none";
        loginForm.style.transform = "translateY(100%)";
        loginForm.removeEventListener("transitionend", handle);
    }

    loginForm.removeEventListener("transitionend", handle);
    loginForm.addEventListener("transitionend", handle);
});

let registrationForm = document.getElementById("registrationForm");
document.getElementById("toLogin").addEventListener("click", () => {
    authCard.classList.remove("activate-registration");
    loginForm.removeAttribute("style");    
    authCard.classList.add("activate-login");

    function handle(e) {
        if (e.propertyName != "transform" || e.target != registrationForm) return;
        registrationForm.style.transition = "none";
        registrationForm.style.transform = "translateY(100%)";
        registrationForm.removeEventListener("transitionend", handle);
    }

    registrationForm.removeEventListener("transitionend", handle);
    registrationForm.addEventListener("transitionend", handle);
});