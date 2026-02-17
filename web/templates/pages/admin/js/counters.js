function openModal(id) {
  document.getElementById(id).classList.remove("hidden");
  document.getElementById(id).classList.add("flex");
}

function closeModal(id) {
  document.getElementById(id).classList.add("hidden");
  document.getElementById(id).classList.remove("flex");
}

async function addCounter(event) {
  event.preventDefault();
  const form = event.target;
  const formData = new FormData(form);

  const data = {
    number: formData.get("number"),
    name: formData.get("name"),
    location: formData.get("location"),
    category_ids: [],
  };

  try {
    const response = await fetch("/admin/api/counters", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });

    if (response.ok) {
      window.location.reload();
    } else {
      const error = await response.json();
      alert(error.error || "Gagal menambahkan loket");
    }
  } catch (error) {
    alert("Network error");
  }
  return false;
}

async function deleteCounter(id) {
  if (!confirm("Apakah Anda yakin ingin menghapus loket ini?")) return;

  try {
    const response = await fetch(`/admin/api/counters/${id}`, {
      method: "DELETE",
    });
    if (response.ok) {
      window.location.reload();
    } else {
      alert("Gagal menghapus loket");
    }
  } catch (error) {
    alert("Network error");
  }
}

async function loadEditCounter(id) {
  try {
    const response = await fetch(`/admin/api/counters/${id}`);
    
    if (!response.ok) {
      alert("Error: " + response.status);
      return;
    }

    const counter = await response.json();
    
    // Handle sql.NullString fields
    let counterName = "";
    let counterLocation = "";
    
    if (counter.name) {
      counterName = typeof counter.name === 'string' ? counter.name : (counter.name.String || "");
    }
    if (counter.location) {
      counterLocation = typeof counter.location === 'string' ? counter.location : (counter.location.String || "");
    }

    document.getElementById("editCounterId").value = counter.id;
    document.getElementById("editCounterNumber").value = counter.number || "";
    document.getElementById("editCounterName").value = counterName;
    document.getElementById("editCounterLocation").value = counterLocation;
    document.getElementById("editCounterTitleName").textContent = counterName || counter.number || "";

    document.getElementById("editCounterModal").classList.remove("hidden");
    document.getElementById("editCounterModal").classList.add("flex");
  } catch (error) {
    console.error(error);
    alert("Error: " + error.message);
  }
}

async function editCounter(event) {
  event.preventDefault();
  const form = event.target;
  const formData = new FormData(form);
  const counterId = formData.get("id");

  const data = {
    number: formData.get("number"),
    name: formData.get("name"),
    location: formData.get("location"),
    category_ids: [],
  };

  try {
    const response = await fetch(`/admin/api/counters/${counterId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });

    if (response.ok) {
      window.location.reload();
    } else {
      const error = await response.json();
      alert(error.error || "Gagal memperbarui loket");
    }
  } catch (error) {
    alert("Network error");
  }
  return false;
}

async function setCounterStatus(id, currentStatus) {
  if (currentStatus === "offline") {
    if (!confirm("Apakah Anda yakin ingin mengaktifkan loket ini?")) return;
  } else {
    if (!confirm("Apakah Anda yakin ingin menonaktifkan (offline) loket ini?")) return;
  }

  try {
    const button = event.target.closest("button");
    const originalIcon = button.innerHTML;
    button.disabled = true;
    button.innerHTML = '<i class="fas fa-spinner fa-spin"></i>';

    const newStatus = currentStatus === "offline" ? "idle" : "offline";

    const response = await fetch(`/admin/api/counters/${id}/status`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ status: newStatus }),
    });

    if (response.ok) {
      const successDiv = document.createElement("div");
      successDiv.className =
        "fixed top-4 right-4 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg z-50";
      successDiv.innerHTML = `<i class="fas fa-check-circle"></i> Loket berhasil ${newStatus === "offline" ? "dinonaktifkan" : "diaktifkan"}!`;
      document.body.appendChild(successDiv);

      setTimeout(() => {
        successDiv.remove();
        window.location.reload();
      }, 1000);
    } else {
      const error = await response.json();
      alert(error.error || "Gagal mengubah status loket");
      button.disabled = false;
      button.innerHTML = originalIcon;
    }
  } catch (error) {
    alert("Terjadi kesalahan jaringan. Mohon periksa koneksi Anda.");
    button.disabled = false;
    button.innerHTML = originalIcon;
  }
}

async function manageCounterCategories(counterId, counterName) {
  try {
    // Load current categories for this counter
    const response = await fetch(`/admin/api/counters/${counterId}/categories`);
    if (!response.ok) {
      alert("Gagal memuat kategori loket");
      return;
    }

    const data = await response.json();
    const assignedCategoryIds = data.category_ids || [];

    // Set counter info
    document.getElementById("manageCategoriesCounterId").value = counterId;
    document.getElementById("manageCategoriesCounterName").textContent = counterName || `Loket ${counterId}`;

    // Reset all checkboxes
    document.querySelectorAll('.category-checkbox').forEach(checkbox => {
      checkbox.checked = false;
    });

    // Check assigned categories
    assignedCategoryIds.forEach(categoryId => {
      const checkbox = document.querySelector(`.category-checkbox[value="${categoryId}"]`);
      if (checkbox) {
        checkbox.checked = true;
      }
    });

    openModal("manageCategoriesModal");
  } catch (error) {
    alert("Terjadi kesalahan jaringan saat memuat kategori");
  }
}

async function saveCounterCategories(event) {
  event.preventDefault();
  const form = event.target;
  const counterId = form.querySelector('[name="counter_id"]').value;

  // Get all checked category IDs
  const checkedBoxes = form.querySelectorAll('input[name="category_ids"]:checked');
  const categoryIds = Array.from(checkedBoxes).map(cb => parseInt(cb.value));

  try {
    const response = await fetch(`/admin/api/counters/${counterId}/categories`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ category_ids: categoryIds }),
    });

    if (response.ok) {
      window.location.reload();
    } else {
      const error = await response.json();
      alert(error.error || "Gagal menyimpan kategori");
    }
  } catch (error) {
    alert("Network error");
  }
  return false;
}
