document.addEventListener("DOMContentLoaded", () => {
  const taskForm = document.getElementById("taskForm");
  const healthStatus = document.getElementById("healthStatus");
  const breakBtn = document.getElementById("breakAppBtn");
  let isBroken = false;

  // Task Form Submission
  if (taskForm) {
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
          titleInput.value = "";
          window.location.reload();
        } else {
          const errorText = await response.text();
          console.error("Server Error Detail:", errorText);
          alert(`Server error: ${response.status}`);
        }
      } catch (err) {
        console.error("Fetch Error:", err);
        alert(
          "Connection failed. Is the backend running at " + backendUrl + "?",
        );
      }
    });
  }

  // Handle "Mark Done" button clicks (Event Delegation scoped inside DOMContentLoaded)
  document.addEventListener("click", async (e) => {
    const markDoneBtn = e.target.closest(".mark-done-btn");
    if (!markDoneBtn) return;

    const taskId = markDoneBtn.getAttribute("data-id");
    if (!taskId) return;

    try {
      const response = await fetch(`/api/tasks/${taskId}/done`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        window.location.reload();
      } else {
        const errorText = await response.text();
        console.error("Failed to mark task done:", errorText);
        alert(`Failed to update task: ${response.status}`);
      }
    } catch (err) {
      console.error("Network Error:", err);
      alert("Connection failed while updating task.");
    }
  });

  // Health status check function
  async function checkHealth() {
    try {
      const response = await fetch("/api/health");
      if (response.ok) {
        isBroken = false;
        if (healthStatus) {
          healthStatus.textContent = "✓ Backend Healthy";
          healthStatus.className = "health-status healthy";
        }
        if (breakBtn) {
          breakBtn.textContent = "Break App";
          breakBtn.className = "break-btn";
        }
      } else {
        handleUnhealthyState();
      }
    } catch (err) {
      console.error("Health check failed:", err);
      handleUnhealthyState();
    }
  }

  function handleUnhealthyState() {
    isBroken = true;
    if (healthStatus) {
      healthStatus.textContent = "✗ Backend Unhealthy";
      healthStatus.className = "health-status unhealthy";
    }
    if (breakBtn) {
      breakBtn.textContent = "Fix App";
      breakBtn.className = "break-btn fix-mode";
    }
  }

  // Break/Fix App button functionality
  if (breakBtn) {
    breakBtn.addEventListener("click", async () => {
      try {
        const endpoint = isBroken ? "/api/fix" : "/api/break";
        const response = await fetch(endpoint, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
        });

        if (response.ok) {
          alert(
            isBroken
              ? "App has been fixed!"
              : "App has been broken! The pod will restart shortly.",
          );
          setTimeout(checkHealth, 1000);
        } else {
          alert("Failed to toggle app state: " + response.status);
        }
      } catch (err) {
        console.error("Toggle Error:", err);
        alert("Failed to toggle app state: " + err.message);
      }
    });
  }

  // Initialize checks
  checkHealth();
  setInterval(checkHealth, 5000);
});
