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

  // Health status indicator
  const healthStatus = document.getElementById("healthStatus");
  const breakBtn = document.getElementById("breakAppBtn");
  let isBroken = false;

  // Function to check health status
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
    } catch (err) {
      console.error("Health check failed:", err);
      isBroken = true;
      if (healthStatus) {
        healthStatus.textContent = "✗ Backend Unreachable";
        healthStatus.className = "health-status unhealthy";
      }
      if (breakBtn) {
        breakBtn.textContent = "Fix App";
        breakBtn.className = "break-btn fix-mode";
      }
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
          if (isBroken) {
            alert("App has been fixed!");
          } else {
            alert("App has been broken! The pod will restart shortly.");
          }
          // Check health status after action
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

  // Initial health check
  checkHealth();

  // Periodic health check every 5 seconds
  setInterval(checkHealth, 5000);
});
