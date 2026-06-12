const taskForm = document.getElementById('taskForm');
const taskInput = document.getElementById('taskInput');
const taskList = document.getElementById('taskList');
const taskCount = document.getElementById('taskCount');

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
    main.textContent = t.text;

    const meta = document.createElement('div');
    meta.className = 'meta';
    meta.textContent = `ID: ${t.id}`;

    left.appendChild(main);
    left.appendChild(meta);

    const pill = document.createElement('div');
    pill.className = `pill ${t.completed ? 'done' : 'pending'}`;
    pill.textContent = t.completed ? 'Completed' : 'Pending';

    li.appendChild(left);
    li.appendChild(pill);

    taskList.appendChild(li);
  }
}

taskForm.addEventListener('submit', async (e) => {
  e.preventDefault();

  const text = taskInput.value.trim();
  if (!text) return;

  const res = await fetch('/api/add', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ text }),
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    alert(err.error || 'Failed to add task');
    return;
  }

  taskInput.value = '';
  await fetchTasks();
});

// Initial load
fetchTasks().catch((err) => {
  console.error(err);
  alert('Could not load tasks from server.');
});

