document.addEventListener('DOMContentLoaded', () => {
    const cart = [];

    document.querySelectorAll('.add-to-cart').forEach(button => {
        button.addEventListener('click', () => {
            const title = button.getAttribute('data-title');
            const price = button.getAttribute('data-price');

            cart.push({ title, price });

            console.log('Book added to cart:', title, price);
            alert(`Book "${title}" has been added to the cart.`);
        });
    });
});

function toggleDescription(button) {
    const bookDetails = button.closest('.book-details');
    const isExpanded = bookDetails.classList.toggle('expanded');
    button.textContent = isExpanded ? 'View Less' : 'View More';
}
