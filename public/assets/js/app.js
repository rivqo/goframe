// Main JavaScript file for GoFrame applications

document.addEventListener("DOMContentLoaded", () => {
  // Initialize components
  initializeNavigation()

  // Add any custom initialization code here
})

function initializeNavigation() {
  // Mobile navigation toggle
  const navToggle = document.querySelector(".nav-toggle")
  const navLinks = document.querySelector(".nav-links")

  if (navToggle && navLinks) {
    navToggle.addEventListener("click", () => {
      navLinks.classList.toggle("active")
    })
  }
}

// Flash message handling
function showFlashMessage(message, type = "info") {
  const flashContainer = document.querySelector(".flash-container")

  if (!flashContainer) {
    const container = document.createElement("div")
    container.className = "flash-container"
    document.body.appendChild(container)
  }

  const flashMessage = document.createElement("div")
  flashMessage.className = `flash-message ${type}`
  flashMessage.textContent = message

  const closeButton = document.createElement("button")
  closeButton.className = "flash-close"
  closeButton.innerHTML = "&times;"
  closeButton.addEventListener("click", () => {
    flashMessage.remove()
  })

  flashMessage.appendChild(closeButton)
  document.querySelector(".flash-container").appendChild(flashMessage)

  // Auto-remove after 5 seconds
  setTimeout(() => {
    if (flashMessage.parentNode) {
      flashMessage.remove()
    }
  }, 5000)
}

// Form validation helper
function validateForm(formSelector, rules) {
  const form = document.querySelector(formSelector)

  if (!form) return

  form.addEventListener("submit", (event) => {
    let isValid = true

    for (const field in rules) {
      const input = form.querySelector(`[name="${field}"]`)
      const errorElement = form.querySelector(`[data-error="${field}"]`)

      if (!input || !errorElement) continue

      const value = input.value.trim()
      let errorMessage = ""

      for (const rule of rules[field]) {
        if (rule.type === "required" && value === "") {
          errorMessage = rule.message || "This field is required"
          isValid = false
          break
        }

        if (rule.type === "email" && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
          errorMessage = rule.message || "Please enter a valid email address"
          isValid = false
          break
        }

        if (rule.type === "minLength" && value.length < rule.value) {
          errorMessage = rule.message || `Must be at least ${rule.value} characters`
          isValid = false
          break
        }
      }

      errorElement.textContent = errorMessage

      if (errorMessage) {
        input.classList.add("is-invalid")
      } else {
        input.classList.remove("is-invalid")
      }
    }

    if (!isValid) {
      event.preventDefault()
    }
  })
}

