function temaOscuro() {
    let body = document.body;
    let icon = document.querySelector('#d1-icon');

    if (body.getAttribute("data-bs-theme") === "light") {
        body.setAttribute("data-bs-theme", "dark");
        icon.className = "bi bi-sun-fill"; // Ícono de sol en modo oscuro
    } else {
        body.setAttribute("data-bs-theme", "light");
        icon.className = "bi bi-moon-fill"; // Ícono de luna en modo claro
    }
}
