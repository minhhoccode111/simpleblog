document.addEventListener("DOMContentLoaded", () => {
  const formDelete = document.getElementById("form-delete-article");
  formDelete.addEventListener("submit", (e) => {
    const confirmed = confirm("Are you sure you want to delete this article?");
    if (!confirmed) {
      e.preventDefault();
    }
  });
});
