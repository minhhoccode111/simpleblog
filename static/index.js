document.addEventListener("DOMContentLoaded", () => {
  const forms = document.querySelectorAll(
    '[data-form-delete-article="form-delete-article"]',
  );
  forms.forEach((form) => {
    form.addEventListener("submit", (e) => {
      const confirmed = confirm(
        "Are you sure you want to delete this article?",
      );
      if (!confirmed) {
        e.preventDefault();
      }
    });
  });
});
