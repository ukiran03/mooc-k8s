document.addEventListener("DOMContentLoaded", () => {
  const taskForm = document.getElementById("taskForm");
  if (!taskForm) return;

  // Pull the URL directly from the template data attribute
  const backendUrl = taskForm.getAttribute("data-url");

  taskForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const titleInput = document.getElementById("taskTitle");
    const data = {
      title: titleInput.value,
      state: 0,
    };

    try {
      const response = await fetch(backendUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      if (response.ok) {
        // Clear input and reload to show the new task
        titleInput.value = "";
        window.location.reload();
      } else {
        const errorText = await response.text();
        console.error("Server Error Detail:", errorText);
        alert(`Server error: ${response.status}`);
      }
    } catch (err) {
      console.error("Fetch Error:", err);
      alert("Connection failed. Is the backend running at " + backendUrl + "?");
    }
  });
});
