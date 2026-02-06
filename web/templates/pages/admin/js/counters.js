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

    const categoryId = formData.get("category_id");
    const data = {
      number: formData.get("number"),
      name: formData.get("name"),
      location: formData.get("location"),
      category_id: categoryId ? parseInt(categoryId) : null,
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
        alert("Gagal memuat data loket");
        return;
      }

      const counter = await response.json();

      document.getElementById("editCounterId").value = counter.id;
      document.getElementById("editCounterNumber").value = counter.number;
      document.getElementById("editCounterName").value = counter.name;
      document.getElementById("editCounterLocation").value =
        counter.location || "";

      const categorySelect = document.getElementById("editCounterCategory");
      categorySelect.value = counter.category_id || "";

      openModal("editCounterModal");
    } catch (error) {
      alert("Terjadi kesalahan jaringan saat memuat data loket");
    }
  }

  async function editCounter(event) {
    event.preventDefault();
    const form = event.target;
    const formData = new FormData(form);
    const counterId = formData.get("id");

    const categoryId = formData.get("category_id");
    const data = {
      number: formData.get("number"),
      name: formData.get("name"),
      location: formData.get("location"),
      category_id: categoryId ? parseInt(categoryId) : null,
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
