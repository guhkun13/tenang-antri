  function staffDashboard() {
    const dataEl = document.getElementById("staff-data");
    const counterNumber = dataEl.dataset.counterNumber || "1";

    function speakTicket(ticketNumber, counterNum) {
      if ('speechSynthesis' in window) {
        window.speechSynthesis.cancel();
        const utterance = new SpeechSynthesisUtterance("Tiket nomor " + ticketNumber + " ke loket " + counterNum);
        utterance.lang = 'id-ID';
        utterance.rate = 0.9;
        utterance.onend = function() {
          console.log('Speech finished');
        };
        window.speechSynthesis.speak(utterance);
      }
    }

    return {
      loading: false,
      hasCurrentTicket: dataEl.dataset.hasTicket === "true",
      counterStatus: dataEl.dataset.counterStatus,
      counterNumber: counterNumber,
      toast: { show: false, message: "", type: "success" },

      showToast: function (message, type) {
        type = type || "success";
        this.toast = { show: true, message: message, type: type };
        var self = this;
        setTimeout(function () {
          self.toast.show = false;
        }, 3000);
      },

      callNext: function () {
        if (this.hasCurrentTicket) {
          this.showToast("Selesaikan tiket saat ini sebelum memanggil berikutnya", "error");
          return;
        }
        console.log("callNext");
        var self = this;
        self.loading = true;
        fetch("/staff/call-next", { method: "POST" })
          .then(function (response) {
            console.log("response", response);
            return response.json();
          })
          .then(function (data) {
            console.log("data", data);
            if (data.error) {
              self.showToast(data.error, "error");
            } else if (data.message) {
              self.showToast(data.message);
            } else {
              self.showToast("Tiket berikutnya dipanggil");
              if (data.ticket_number) {
                speakTicket(data.ticket_number, self.counterNumber);
              }
              self.hasCurrentTicket = true;
              setTimeout(function () {
                window.location.reload();
              }, 500);
            }
          })
          .catch(function (error) {
            console.log("error", error);
            self.showToast("Network error", "error");
          })
          .finally(function () {
            console.log("finally");
            self.loading = false;
          });
      },

      callAgain: function () {
        var self = this;
        if (!this.hasCurrentTicket) {
          this.showToast("Tidak ada ticket yang sedang dilayani", "error");
          return;
        }
        self.loading = true;
        fetch("/staff/call-again", { method: "POST" })
          .then(function (response) {
            return response.json();
          })
          .then(function (data) {
            if (data.error) {
              self.showToast(data.error, "error");
            } else {
              self.showToast("Tiket dipanggil ulang");
              if (data.ticket_number) {
                speakTicket(data.ticket_number, self.counterNumber);
              }
              setTimeout(function () {
                window.location.reload();
              }, 500);
            }
          })
          .catch(function (error) {
            self.showToast("Network error", "error");
          })
          .finally(function () {
            self.loading = false;
          });
      },

      completeTicket: function () {
        var self = this;
        self.loading = true;
        fetch("/staff/complete", { method: "POST" })
          .then(function (response) {
            return response.json();
          })
          .then(function (data) {
            if (data.error) {
              self.showToast(data.error, "error");
            } else {
              self.showToast("Tiket selesai");
              self.hasCurrentTicket = false;
              setTimeout(function () {
                window.location.reload();
              }, 500);
            }
          })
          .catch(function (error) {
            self.showToast("Network error", "error");
          })
          .finally(function () {
            self.loading = false;
          });
      },

      markNoShow: function () {
        var self = this;
        self.loading = true;
        fetch("/staff/no-show", { method: "POST" })
          .then(function (response) {
            return response.json();
          })
          .then(function (data) {
            if (data.error) {
              self.showToast(data.error, "error");
            } else {
              self.showToast("Ditandai tidak hadir");
              self.hasCurrentTicket = false;
              setTimeout(function () {
                window.location.reload();
              }, 500);
            }
          })
          .catch(function (error) {
            self.showToast("Network error", "error");
          })
          .finally(function () {
            self.loading = false;
          });
      },

      toggleCounterStatus: function () {
        var self = this;
        if (self.counterStatus === "disabled") {
          self.showToast("Loket dinonaktifkan oleh admin", "error");
          return;
        }
        var action = self.counterStatus === "idle" ? "pause" : "resume";
        fetch("/staff/" + action, { method: "POST" })
          .then(function (response) {
            return response.json();
          })
          .then(function (data) {
            if (data.error) {
              self.showToast(data.error, "error");
            } else {
              self.counterStatus = action === "pause" ? "paused" : "idle";
              self.showToast("Loket " + (action === "pause" ? "dijeda" : "dilanjutkan"));
            }
          })
          .catch(function (error) {
            self.showToast("Network error", "error");
          });
      },
    };
  }

  function headerStatus() {
    const dataEl = document.getElementById("staff-data");
    return {
      counterStatus: dataEl ? dataEl.dataset.counterStatus : 'idle'
    };
  }