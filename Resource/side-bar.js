document.querySelectorAll('.menu-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        const sideMenu = document.getElementById('side-menu');
        const isClosing = sideMenu.classList.toggle('closed');

        // Toggle menu icon rotation
        document.querySelectorAll('.hamburger-menu-icon').forEach(icon => {
            icon.style.transform = isClosing ? 'rotate(0deg)' : 'rotate(90deg)';
        });
    });
});