<template>
  <div class="mb-4 p-4 rounded-lg">
    <h3 class="text-xl font-semibold mb-3">Add New Task</h3>
    <form @submit.prevent="submitTask" class="flex space-x-2">
      <input 
        v-model="newTask" 
        type="text" 
        placeholder="Enter a new task" 
        class="flex-grow px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        required
      />
      <button 
        type="submit" 
        class="bg-blue text-white px-4 py-2 rounded-md hover:bg-blue-600 transition-colors"
      >
        Add Task
      </button>
    </form>
    <p v-if="errorMessage" class="text-red-500 mt-2">{{ errorMessage }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue';

// Emit an event to parent component
const emit = defineEmits(['task-added']);

// State for new task input
const newTask = ref('');
const errorMessage = ref('');

// Submit task method
const submitTask = async () => {
  // Trim the task and validate
  const taskText = newTask.value.trim();
  if (!taskText) {
    errorMessage.value = 'Task cannot be empty';
    return;
  }

  try {
    // Send POST request to create task
    const response = await fetch('http://localhost:8080/tasks', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        task: taskText,
        completed: false
      })
    });

    if (!response.ok) {
      throw new Error('Failed to create task');
    }

    // Emit event to parent component
    emit('task-added');

    // Reset input and clear any previous error
    newTask.value = '';
    errorMessage.value = '';
  } catch (error) {
    console.error('Error creating task:', error);
    errorMessage.value = 'Failed to add task. Please try again.';
  }
};
</script>