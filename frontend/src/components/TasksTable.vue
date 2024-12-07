<script setup>
import { ref, onMounted, onUnmounted } from 'vue';

const tasks = ref([]);
const message = ref('');
const ws = ref(null);

// Complete a task
const completeTask = async (taskId) => {
  try {
    const response = await fetch(`http://localhost:8080/tasks/${taskId}/complete`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error('Failed to complete task');
    }

    // The WebSocket will handle updating the tasks list
  } catch (error) {
    console.error('Error completing task:', error);
    message.value = 'Failed to complete task';
  }
};

// Delete a task
const deleteTask = async (taskId) => {
  try {
    const response = await fetch(`http://localhost:8080/tasks/${taskId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error('Failed to delete task');
    }

    // The WebSocket will handle updating the tasks list
  } catch (error) {
    console.error('Error deleting task:', error);
    message.value = 'Failed to delete task';
  }
};

// Fetch the initial tasks
const fetchTasks = async () => {
  try {
    const response = await fetch('http://localhost:8080/tasks');
    if (response.ok) {
      const fetchedTasks = await response.json();
      
      tasks.value = Array.isArray(fetchedTasks) ? fetchedTasks : [];
      
      if (tasks.value.length === 0) {
        message.value = 'No tasks available';
      } else {
        message.value = '';
      }
    } else {
      message.value = 'Failed to fetch tasks';
      console.error('Failed to fetch tasks');
    }
  } catch (error) {
    message.value = 'Error fetching tasks';
    console.error('Error fetching tasks:', error);
  }
};

// Setup WebSocket connection
const setupWebSocket = () => {
  try {
    ws.value = new WebSocket('ws://localhost:8080/tasks/ws');

    ws.value.onopen = () => {
      console.log('WebSocket connection established');
      message.value = '';
    };

    ws.value.onmessage = (event) => {
      try {
        const updatedTasks = JSON.parse(event.data);
        
        if (Array.isArray(updatedTasks)) {
          tasks.value = updatedTasks;
          message.value = updatedTasks.length === 0 ? 'No tasks available' : '';
        } else {
          tasks.value = [];
          message.value = 'Invalid tasks data received';
        }
      } catch (parseError) {
        console.error('Error parsing WebSocket message:', parseError);
        message.value = 'Error processing task updates';
      }
    };

    ws.value.onerror = (error) => {
      console.error('WebSocket error:', error);
      message.value = 'WebSocket connection error';
    };

    ws.value.onclose = (event) => {
      console.log('WebSocket connection closed');
      if (!event.wasClean) {
        message.value = 'WebSocket connection unexpectedly closed';
      }
    };
  } catch (error) {
    console.error('Error setting up WebSocket:', error);
    message.value = 'Failed to establish WebSocket connection';
  }
};

// Cleanup WebSocket on component unmount
const cleanupWebSocket = () => {
  if (ws.value) {
    ws.value.close();
    ws.value = null;
  }
};

// Fetch tasks and setup WebSocket on mount
onMounted(() => {
  fetchTasks();
  setupWebSocket();
});

// Cleanup WebSocket when component is unmounted
onUnmounted(cleanupWebSocket);
</script>

<template>
  <div class="w-full p-4 bg-white rounded-lg shadow-md">
    <h2 class="text-2xl font-bold mb-4">Task List</h2>

    <!-- Display error or no tasks message -->
    <p v-if="message" class="text-red-500">{{ message }}</p>

    <!-- Table to display tasks -->
    <table v-if="tasks.length > 0" class="w-full text-left table-auto">
      <thead>
        <tr>
          <th class="px-4 py-2 border-b">Task</th>
          <th class="px-4 py-2 border-b">Created At</th>
          <th class="px-4 py-2 border-b">Completed</th>
          <th class="px-4 py-2 border-b">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="task in tasks" :key="task.id">
          <td class="px-4 py-2 border-b">{{ task.task || 'Unnamed Task' }}</td>
          <td class="px-4 py-2 border-b">{{ task.createdAt || 'N/A' }}</td>
          <td class="px-4 py-2 border-b">{{ task.completed ? 'Yes' : 'No' }}</td>
          <td class="px-4 py-2 border-b">
            <button 
              v-if="!task.completed" 
              @click="completeTask(task.id)" 
              class="bg-green-500 text-white px-2 py-1 rounded hover:bg-green-600"
            >
              Complete
            </button>
            <button 
              @click="deleteTask(task.id)" 
              class="bg-red-500 text-white px-2 py-1 rounded hover:bg-red-600 ml-2"
            >
              Delete
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
