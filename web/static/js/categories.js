let selectedIcon = "users";
let editSelectedIcon = "users";

function openModal(id) {
  const modal = document.getElementById(id);
  modal.classList.remove("hidden");
  modal.classList.add("flex");

  document.addEventListener("keydown", handleEscapeKey);

  setTimeout(() => {
    const firstInput = modal.querySelector('input[type="text"]');
    if (firstInput) {
      firstInput.focus();
      firstInput.select();
    }
  }, 100);
}

function closeModal(id) {
  const modal = document.getElementById(id);
  modal.classList.add("hidden");
  modal.classList.remove("flex");

  document.removeEventListener("keydown", handleEscapeKey);
}

function handleEscapeKey(event) {
  if (event.key === "Escape") {
    const modals = document.querySelectorAll('[id$="Modal"]');
    modals.forEach((modal) => {
      if (!modal.classList.contains("hidden")) {
        closeModal(modal.id);
      }
    });
  }
}

function selectIcon(icon) {
  selectedIcon = icon;
  document.getElementById("selectedIcon").value = icon;

  document.querySelectorAll(".icon-option").forEach((btn) => {
    btn.classList.remove("bg-blue-100", "border-blue-500");
  });
  document
    .querySelector(`[data-icon="${icon}"]`)
    .classList.add("bg-blue-100", "border-blue-500");
}

function selectEditIcon(icon) {
  editSelectedIcon = icon;
  document.getElementById("editSelectedIcon").value = icon;

  document.querySelectorAll(".edit-icon-option").forEach((btn) => {
    btn.classList.remove("bg-blue-100", "border-blue-500");
  });
  document
    .querySelector(`.edit-icon-option[data-icon="${icon}"]`)
    .classList.add("bg-blue-100", "border-blue-500");
}

function validateName(input) {
  const value = input.value.trim();
  const errorSpan = document.getElementById(input.id + "Error");

  if (value.length < 2) {
    input.setCustomValidity("Name must be at least 2 characters");
    if (errorSpan) {
      errorSpan.textContent = "Name must be at least 2 characters";
      errorSpan.classList.remove("hidden");
    }
    input.classList.add("border-red-500");
    return false;
  }
  if (value.length > 100) {
    input.setCustomValidity("Name must be less than 100 characters");
    if (errorSpan) {
      errorSpan.textContent = "Name must be less than 100 characters";
      errorSpan.classList.remove("hidden");
    }
    input.classList.add("border-red-500");
    return false;
  }
  input.setCustomValidity("");
  if (errorSpan) {
    errorSpan.classList.add("hidden");
  }
  input.classList.remove("border-red-500");
  return true;
}

function validatePrefix(input) {
  const value = input.value.trim();
  const errorSpan = document.getElementById(input.id + "Error");

  if (value.length < 1) {
    input.setCustomValidity("Prefix is required");
    if (errorSpan) {
      errorSpan.textContent = "Prefix is required";
      errorSpan.classList.remove("hidden");
    }
    input.classList.add("border-red-500");
    return false;
  }
  if (value.length > 10) {
    input.setCustomValidity("Prefix must be less than 10 characters");
    if (errorSpan) {
      errorSpan.textContent = "Prefix must be less than 10 characters";
      errorSpan.classList.remove("hidden");
    }
    input.classList.add("border-red-500");
    return false;
  }
  if (!/^[A-Za-z0-9]+$/.test(value)) {
    input.setCustomValidity("Prefix must contain only letters and numbers");
    if (errorSpan) {
      errorSpan.textContent = "Prefix must contain only letters and numbers";
      errorSpan.classList.remove("hidden");
    }
    input.classList.add("border-red-500");
    return false;
  }
  input.setCustomValidity("");
  if (errorSpan) {
    errorSpan.classList.add("hidden");
  }
  input.classList.remove("border-red-500");
  return true;
}

function validatePriority(input) {
  const value = parseInt(input.value);
  const errorSpan = document.getElementById(input.id + "Error");

  if (isNaN(value) || value < 0 || value > 100) {
    input.setCustomValidity("Priority must be a number between 0 and 100");
    if (errorSpan) {
      errorSpan.textContent = "Priority must be a number between 0 and 100";
      errorSpan.classList.remove("hidden");
    }
    input.classList.add("border-red-500");
    return false;
  }
  input.setCustomValidity("");
  if (errorSpan) {
    errorSpan.classList.add("hidden");
  }
  input.classList.remove("border-red-500");
  return true;
}

async function addCategory(event) {
  event.preventDefault();
  const form = event.target;
  const formData = Object.fromEntries(new FormData(form));
  
  const data = {
    ...formData,
    priority: parseInt(formData.priority) || 0
  };

  console.log('Category data being sent:', data);

  try {
    const response = await fetch("/admin/api/categories", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });

    if (response.ok) {
      const successDiv = document.createElement('div');
      successDiv.className = 'fixed top-4 right-4 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg z-50';
      successDiv.innerHTML = '<i class="fas fa-check-circle"></i> Category added successfully!';
      document.body.appendChild(successDiv);
      
      setTimeout(() => {
        successDiv.remove();
        window.location.reload();
      }, 1500);
    } else {
      const error = await response.json();
      alert(error.error || "Failed to add category");
    }
  } catch (error) {
    alert("Network error. Please check your connection.");
  }
  return false;
}

async function loadCategoryForEdit(id) {
  try {
    const loadingDiv = document.createElement("div");
    loadingDiv.className =
      "fixed top-4 right-4 bg-blue-500 text-white px-6 py-3 rounded-lg shadow-lg z-50";
    loadingDiv.innerHTML =
      '<i class="fas fa-spinner fa-spin"></i> Loading category data...';
    document.body.appendChild(loadingDiv);

    const response = await fetch(`/admin/api/categories/${id}`);
    loadingDiv.remove();

    if (response.ok) {
      const category = await response.json();

      console.log("Category data loaded:", category);

      const form = document.getElementById("editCategoryForm");
      form.reset();

      document.querySelectorAll('[id$="Error"]').forEach((span) => {
        span.classList.add("hidden");
      });
      document.querySelectorAll(".border-red-500").forEach((input) => {
        input.classList.remove("border-red-500");
      });

      document.getElementById("editCategoryId").value = category.id || "";
      document.getElementById("editName").value = category.name || "";
      document.getElementById("editPrefix").value = category.prefix || "";
      document.getElementById("editPriority").value = category.priority || 0;
      document.getElementById("editColorCode").value =
        category.color_code || "#3B82F6";
      document.getElementById("editDescription").value =
        category.description || "";

      setTimeout(() => {
        selectEditIcon(category.icon || "users");
      }, 100);

      openModal("editCategoryModal");
    } else {
      const error = await response.json();
      console.error("Failed to load category:", error);
      alert(error.error || "Failed to load category data");
    }
  } catch (error) {
    console.error("Network error in loadCategoryForEdit:", error);
    alert("Network error. Please check your connection.");
  }
}

async function updateCategory(event) {
  event.preventDefault();
  const form = event.target;
  const id = form.id.value;
  const formData = Object.fromEntries(new FormData(form));
  delete formData.id;

  const data = {
    ...formData,
    priority: parseInt(formData.priority) || 0
  };

  console.log('Category update data being sent:', data);

  if (!data.name || !data.prefix || !data.icon) {
    alert("Please fill in all required fields");
    return false;
  }

  if (!/^[A-Za-z0-9]+$/.test(data.prefix)) {
    alert("Prefix should contain only letters and numbers");
    return false;
  }

  if (isNaN(data.priority) || data.priority < 0 || data.priority > 100) {
    alert("Priority should be a number between 0 and 100");
    return false;
  }

  if (!/^#[0-9A-F]{6}$/i.test(data.color_code)) {
    alert("Please enter a valid color code (e.g., #3B82F6)");
    return false;
  }

  try {
    const submitButton = form.querySelector('button[type="submit"]');
    submitButton.disabled = true;
    submitButton.innerHTML =
      '<i class="fas fa-spinner fa-spin"></i> Updating...';

    const response = await fetch(`/admin/api/categories/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });

    if (response.ok) {
      const successDiv = document.createElement("div");
      successDiv.className =
        "fixed top-4 right-4 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg z-50";
      successDiv.innerHTML =
        '<i class="fas fa-check-circle"></i> Category updated successfully!';
      document.body.appendChild(successDiv);

      setTimeout(() => {
        successDiv.remove();
        window.location.reload();
      }, 1500);
    } else {
      const error = await response.json();
      alert(error.error || "Failed to update category");
    }
  } catch (error) {
    alert("Network error. Please check your connection.");
  } finally {
    const submitButton = form.querySelector('button[type="submit"]');
    submitButton.disabled = false;
    submitButton.innerHTML = "Update Category";
  }
  return false;
}

async function toggleCategoryStatus(id, currentStatus) {
  const action = currentStatus === "active" ? "deactivate" : "activate";
  if (!confirm(`Are you sure you want to ${action} this category?`)) return;

  try {
    const newStatus = currentStatus === "active" ? false : true;

    const button = event.target.closest("button");
    const originalIcon = button.innerHTML;
    button.disabled = true;
    button.innerHTML = '<i class="fas fa-spinner fa-spin"></i>';

    const response = await fetch(`/admin/api/categories/${id}/status`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ is_active: newStatus }),
    });

    if (response.ok) {
      const successDiv = document.createElement("div");
      successDiv.className =
        "fixed top-4 right-4 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg z-50";
      successDiv.innerHTML = `<i class="fas fa-check-circle"></i> Category ${action}d successfully!`;
      document.body.appendChild(successDiv);

      setTimeout(() => {
        successDiv.remove();
        window.location.reload();
      }, 1000);
    } else {
      const error = await response.json();
      alert(error.error || `Failed to ${action} category`);
      button.disabled = false;
      button.innerHTML = originalIcon;
    }
  } catch (error) {
    alert("Network error. Please check your connection.");
    button.disabled = false;
    button.innerHTML = originalIcon;
  }
}

async function deleteCategory(id) {
  if (
    !confirm(
      "Are you sure you want to delete this category? This action cannot be undone.",
    )
  )
    return;

  try {
    const button = event.target.closest("button");
    const originalIcon = button.innerHTML;
    button.disabled = true;
    button.innerHTML = '<i class="fas fa-spinner fa-spin"></i>';

    const response = await fetch(`/admin/api/categories/${id}`, {
      method: "DELETE",
    });

    if (response.ok) {
      const successDiv = document.createElement("div");
      successDiv.className =
        "fixed top-4 right-4 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg z-50";
      successDiv.innerHTML =
        '<i class="fas fa-check-circle"></i> Category deleted successfully!';
      document.body.appendChild(successDiv);

      setTimeout(() => {
        successDiv.remove();
        window.location.reload();
      }, 1000);
    } else {
      const error = await response.json();
      alert(error.error || "Failed to delete category");
      button.disabled = false;
      button.innerHTML = originalIcon;
    }
  } catch (error) {
    alert("Network error. Please check your connection.");
    button.disabled = false;
    button.innerHTML = originalIcon;
  }
}

function editCategory(id) {
  loadCategoryForEdit(id);
}

function filterCategories() {
  const searchTerm = document
    .getElementById("categorySearch")
    .value.toLowerCase();
  const statusFilter = document.getElementById("filterStatus").value;
  const cards = document.querySelectorAll(".category-card");

  cards.forEach((card) => {
    const name = card.dataset.name.toLowerCase();
    const status = card.dataset.status;
    const matchesSearch = name.includes(searchTerm);
    const matchesStatus = !statusFilter || status === statusFilter;

    if (matchesSearch && matchesStatus) {
      card.style.display = "";
    } else {
      card.style.display = "none";
    }
  });
}

document.addEventListener("DOMContentLoaded", function () {
  selectIcon("users");
  selectEditIcon("users");
});
