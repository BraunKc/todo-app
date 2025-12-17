const API_ADDR = "http://localhost:8080";

function saveTheme(theme) {
    localStorage.setItem("theme", theme);
}

function setTheme() {
    document.body.classList.toggle("light");
    const isLight = document.getElementById("theme-btn").classList.toggle("light");
    saveTheme(isLight);
}

document.getElementById("theme-btn").addEventListener("click", setTheme);

let title = "";
let orderByField = 0;
let orderByDirection = 0;
let taskStatuses = [0, 1, 2];
let taskPriorities = [0, 1, 2];

function filters() {
    const isFiltersOpen = document.getElementById("filters").classList.toggle("show");

    if (isFiltersOpen) return;

    const checkedStatuses = document.querySelectorAll("#filters-statuses input.filters-checkbox:checked");
    const checkedPriorities = document.querySelectorAll("#filters-priorities input.filters-checkbox:checked");

    taskStatuses = checkedStatuses.length ? Array.from(checkedStatuses).map(status => Number(status.value)) : [0, 1, 2];

    taskPriorities = checkedPriorities.length ? Array.from(checkedPriorities).map(priority => Number(priority.value)) : [0, 1, 2];

    loadTasks(taskStatuses, taskPriorities, orderByField, orderByDirection, title);
}

document.getElementById("filters-btn")?.addEventListener("click", filters);

// textarea autoresize
function autoResize(e) {
    e.target.style.height = "auto";
    e.target.style.height = e.target.scrollHeight + "px";
}

// CRUD
async function createTask(title, description, priority, dueDate) {
    try {
        const resp = await fetch(`${API_ADDR}/api/v1/task/`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title, description, priority, due_date: dueDate })
        });

        if (!resp.ok) throw new Error(`HTTP error. Status: ${resp.status}`);
        const data = await resp.json();
        return data.task;
    } catch (error) {
        console.error("Failed to create task:", error);
    }
}

async function updateTask(id, title, description, status, priority, dueDate) {
    try {
        const resp = await fetch(`${API_ADDR}/api/v1/task/`, {
            method: "PATCH",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ id, title, description, status, priority, due_date: dueDate })
        });

        if (!resp.ok) throw new Error(`HTTP error. Status: ${resp.status}`);
        const data = await resp.json();
        console.log(data);
    } catch (error) {
        console.error("Failed to update task:", error);
    }
}

async function deleteTasks(ids) {
    try {
        const resp = await fetch(`${API_ADDR}/api/v1/task/`, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ ids })
        });

        if (!resp.ok) throw new Error(`HTTP error. Status: ${resp.status}`);
    } catch (error) {
        console.error("Failed to delete tasks:", error);
    }
}

function createSelect(className, options, selectedValue) {
    const select = document.createElement("select");
    select.classList.add(className);

    options.forEach(({ text, value }) => {
        const option = document.createElement("option");
        option.innerText = text;
        option.value = value;
        if (value == selectedValue) option.selected = true;
        select.appendChild(option);
    });

    return select;
}

async function taskEvent(e) {
    const task = e.target.closest(".task");
    if (!task) return;

    if (e.target.classList.contains("delete-btn") && task.id) {
        deleteTasks([task.id]);
        return;
    }

    const titleEl = task.querySelector(".task-title");
    const descriptionEl = task.querySelector(".task-description");
    const priorityEl = task.querySelector(".task-priority");
    const dueDateEl = task.querySelector(".task-due-date");

    const titleValue = titleEl.value;
    const descriptionValue = descriptionEl.value;
    const priorityValue = Number(priorityEl.value);
    const dueDateValue = dueDateEl.value ? new Date(dueDateEl.value).getTime() / 1000 : 0;

    if (!task.id) {
        if (titleValue && descriptionValue && dueDateValue > 0) {
            const newTask = await createTask(titleValue, descriptionValue, priorityValue, dueDateValue);
            if (!newTask) return;

            task.id = newTask.id;

            const statusSelect = createSelect("task-status", [
                { text: "status todo", value: 0 },
                { text: "status in progress", value: 1 },
                { text: "status done", value: 2 }
            ], newTask.status);

            const taskOptions = task.querySelector(".task-options");
            taskOptions.insertBefore(statusSelect, taskOptions.lastElementChild);
        }
    } else if (titleValue && descriptionValue && priorityValue && dueDateValue > 0) {
        const statusValue = Number(task.querySelector(".task-status").value || 0);
        updateTask(task.id, titleValue, descriptionValue, statusValue, priorityValue, dueDateValue);
    }
}

document.getElementById("search").addEventListener("input", (e) => {
    title = e.target.value;
    loadTasks(taskStatuses, taskPriorities, orderByField, orderByDirection, title);
});

document.getElementById("order-by").addEventListener("change", (e) => {
    const mapping = {
        "0": [0, 0],
        "1": [0, 1],
        "2": [1, 0],
        "3": [1, 1],
        "4": [2, 0],
        "5": [2, 1]
    };

    [orderByField, orderByDirection] = mapping[e.target.value];
    loadTasks(taskStatuses, taskPriorities, orderByField, orderByDirection, title);
});

document.getElementById("create-btn").addEventListener("click", () => {
    const task = document.createElement("div");
    task.classList.add("task");

    const titleInput = document.createElement("input");
    titleInput.type = "text";
    titleInput.classList.add("task-title");
    titleInput.placeholder = "title...";

    const descriptionTextarea = document.createElement("textarea");
    descriptionTextarea.classList.add("task-description");
    descriptionTextarea.rows = 1;
    descriptionTextarea.placeholder = "description...";
    descriptionTextarea.addEventListener("input", autoResize);

    const taskOptions = document.createElement("div");
    taskOptions.classList.add("task-options");

    const prioritySelect = createSelect("task-priority", [
        { text: "priority low", value: 0 },
        { text: "priority medium", value: 1 },
        { text: "priority high", value: 2 }
    ], 0);

    const dueDateInput = document.createElement("input");
    dueDateInput.type = "date";
    dueDateInput.classList.add("task-due-date");

    taskOptions.appendChild(prioritySelect);
    taskOptions.appendChild(dueDateInput);

    const deleteBtn = document.createElement("button");
    deleteBtn.classList.add("delete-btn", "material-symbols-outlined");
    deleteBtn.innerText = "delete";

    task.appendChild(titleInput);
    task.appendChild(descriptionTextarea);
    task.appendChild(taskOptions);
    task.appendChild(deleteBtn);

    task.addEventListener("focusout", taskEvent);
    document.getElementById("container").appendChild(task);

    task.scrollIntoView({ behavior: "smooth" });
});

document.getElementById("container").addEventListener("click", (e) => {
    if (e.target.closest(".delete-btn")) {
        e.target.closest(".task").remove();
    }
});

document.getElementById("logout-btn").addEventListener("click", async () => {
    try {
        const resp = await fetch(`${API_ADDR}/api/v1/logout`, {
            method: "POST",
            headers: { "Content-Type": "application/json" }
        });

        if (!resp.ok) throw new Error(`HTTP error. Status: ${resp.status}`);
        window.location.href = "/";
    } catch (error) {
        console.error("Logout failed:", error);
    }
});

async function loadTasks(statuses, priorities, field, direction, searchTitle) {
    try {
        const resp = await fetch(`${API_ADDR}/api/v1/tasks`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                page_size: 1000,
                page_number: 1,
                filters: {
                    task_statuses: statuses,
                    task_priorities: priorities
                },
                order_by: { field, direction },
                title: searchTitle
            })
        });

        if (!resp.ok) throw new Error(`HTTP error. Status: ${resp.status}`);
        const data = await resp.json();

        const container = document.getElementById("container");
        if (!container) return;

        container.querySelectorAll(".task").forEach(el => el.remove());

        data.tasks.forEach(t => {
            const task = document.createElement("div");
            task.classList.add("task");
            task.id = t.id;

            const titleInput = document.createElement("input");
            titleInput.type = "text";
            titleInput.classList.add("task-title");
            titleInput.value = t.title;
            titleInput.placeholder = "title...";

            const descriptionTextarea = document.createElement("textarea");
            descriptionTextarea.classList.add("task-description");
            descriptionTextarea.rows = 1;
            descriptionTextarea.value = t.description;
            descriptionTextarea.placeholder = "description...";
            descriptionTextarea.addEventListener("input", autoResize);

            const taskOptions = document.createElement("div");
            taskOptions.classList.add("task-options");

            const prioritySelect = createSelect("task-priority", [
                { text: "priority low", value: 0 },
                { text: "priority medium", value: 1 },
                { text: "priority high", value: 2 }
            ], t.priority);

            const statusSelect = createSelect("task-status", [
                { text: "status todo", value: 0 },
                { text: "status in progress", value: 1 },
                { text: "status done", value: 2 }
            ], t.status);

            const dueDateInput = document.createElement("input");
            dueDateInput.type = "date";
            dueDateInput.classList.add("task-due-date");

            const dueDate = new Date(t.due_date * 1000);
            dueDateInput.value = dueDate.toISOString().split("T")[0];

            taskOptions.appendChild(prioritySelect);
            taskOptions.appendChild(statusSelect);
            taskOptions.appendChild(dueDateInput);

            const deleteBtn = document.createElement("button");
            deleteBtn.classList.add("delete-btn", "material-symbols-outlined");
            deleteBtn.innerText = "delete";

            task.appendChild(titleInput);
            task.appendChild(descriptionTextarea);
            task.appendChild(taskOptions);
            task.appendChild(deleteBtn);

            task.addEventListener("focusout", taskEvent);
            container.appendChild(task);
        });
    } catch (error) {
        console.error("Failed to load tasks:", error);
    }
}

document.addEventListener("DOMContentLoaded", loadTasks([0, 1, 2], [0, 1, 2], 0, 0, ""));