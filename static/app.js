const taskForm = document.getElementById('taskForm');
const taskInput = document.getElementById('taskInput');
const taskList = document.getElementById('taskList');
const taskCount = document.getElementById('taskCount');

const viewIdInput = document.getElementById('viewIdInput');
const viewTaskBtn = document.getElementById('viewTaskBtn');
const viewResult = document.getElementById('viewResult');

function showError(msg) {
  viewResult.textContent = msg;
  viewResult.classList.remove('ok');
  viewResult.classList.add('err');
}

function showOk(msg) {
  viewResult.textContent = msg;
  viewResult.classList.remove('err');
  viewResult.classList.add('ok');
}


async function fetchTasks() {
  const res = await fetch('/api/tasks');
  if (!res.ok) throw new Error('Failed to load tasks');
  const data = await res.json();
  const tasks = data.tasks || [];

  taskList.innerHTML = '';
  taskCount.textContent = String(tasks.length);

  for (const t of tasks) {
    const li = document.createElement('li');
    li.className = 'task-item';

    const left = document.createElement('div');
    left.className = 'task-text';

    const main = document.createElement('div');
    main.className = 'main';
    main.textContent = t.title;

    const meta = document.createElement('div');
    meta.className = 'meta';
    meta.textContent = `ID: ${t.id}`;

    left.appendChild(main);
    left.appendChild(meta);

    const pill = document.createElement('div');
    pill.className = `pill ${t.done ? 'done' : 'pending'}`;
    pill.textContent = t.done ? 'Completed' : 'Pending';

    const actions = document.createElement('div');
    actions.className = 'actions';

    const toggleBtn = document.createElement('button');
    toggleBtn.className = 'action-btn';
    toggleBtn.type = 'button';
    toggleBtn.textContent = t.done ? 'Mark Pending' : 'Mark Done';
    toggleBtn.addEventListener('click', async () => {
      await fetch(`/api/done?id=${encodeURIComponent(t.id)}`, { method: 'PUT' });
      await fetchTasks();
    });

    const delBtn = document.createElement('button');
    delBtn.className = 'action-btn danger';
    delBtn.type = 'button';
    delBtn.textContent = 'Delete';
    delBtn.addEventListener('click', async () => {
      await fetch(`/api/delete?id=${encodeURIComponent(t.id)}`, { method: 'DELETE' });
      await fetchTasks();
    });

    actions.appendChild(toggleBtn);
    actions.appendChild(delBtn);

    li.appendChild(left);
    li.appendChild(pill);
    li.appendChild(actions);

    taskList.appendChild(li);
  }
}


taskForm.addEventListener('submit', async (e) => {
  e.preventDefault();

  const title = taskInput.value.trim();
  if (!title) return;

  const res = await fetch('/api/add', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title }),
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    alert(err.error || 'Failed to add task');
    return;
  }

  taskInput.value = '';
  await fetchTasks();
});

viewTaskBtn.addEventListener('click', async () => {
  const id = viewIdInput.value.trim();
  if (!id) {
    showError('Enter a task ID');
    return;
  }

  try {
    const res = await fetch(`/api/task?id=${encodeURIComponent(id)}`);
    const data = await res.json().catch(() => ({}));

    if (!res.ok) {
      showError(data.error || 'Task not found');
      return;
    }

    const t = data.task;
    showOk(`Task #${t.id}: ${t.title} | ${t.done ? 'Completed' : 'Pending'}`);
  } catch (e) {
    showError('Request failed');
  }
});

// Initial load
fetchTasks().catch((err) => {
  console.error(err);
  alert('Could not load tasks from server.');
});


