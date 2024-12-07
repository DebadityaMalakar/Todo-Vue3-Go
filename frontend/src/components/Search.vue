<script>
export default {
  name: "Search",
  data() {
    return {
      searchQuery: "",  // The query entered by the user
      tasks: [],        // To store all fetched tasks
      filteredTasks: [] // To store filtered tasks based on search
    };
  },
  watch: {
    // Watch for changes to the search query
    searchQuery(newQuery) {
      this.filterTasks(newQuery);  // Filter tasks whenever the search query changes
    }
  },
  mounted() {
    this.fetchTasks();  // Fetch tasks when the component is mounted
  },
  methods: {
    async fetchTasks() {
      try {
        const response = await fetch('http://localhost:8080/tasks');
        if (response.ok) {
          const data = await response.json();
          this.tasks = data;  // Store the tasks received from the API
          this.filteredTasks = data; // Initially, all tasks are shown
        } else {
          console.error("Failed to fetch tasks:", response.statusText);
        }
      } catch (error) {
        console.error("Error fetching tasks:", error);
      }
    },

    // Method to filter tasks based on the search query
    filterTasks(query) {
      if (query) {
        const regex = new RegExp(query, "i"); // Create a case-insensitive regex
        this.filteredTasks = this.tasks.filter(task => regex.test(task.task)); // Filter tasks by task name
      } else {
        this.filteredTasks = this.tasks;  // If no query, show all tasks
      }
    }
  }
};
</script>

<template>
  <div class="flex justify-center items-center w-full p-4">
    <div class="relative w-full max-w-md group"> <!-- Added 'group' class here -->
      <!-- Search Input -->
      <input 
        type="text" 
        placeholder="Search..." 
        v-model="searchQuery" 
        class="w-full p-3 pl-10 text-sm border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
      />
      <!-- Search Icon -->
      <div class="absolute inset-y-0 left-3 flex items-center text-gray-400">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zm-7 4a7 7 0 1113.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.817-4.817A7 7 0 011 8z" clip-rule="evenodd" />
        </svg>
      </div>

      <!-- Hoverable Task List -->
      <div 
        v-if="searchQuery" 
        class="mt-4 w-full max-w-md absolute left-0 bg-white shadow-lg rounded-lg border border-gray-200 group-hover:block hidden"
      >
        <ul>
          <li v-for="task in filteredTasks" :key="task.createdAt" class="p-2 border-b border-gray-300">
            {{ task.task }}
          </li>
          <li v-if="filteredTasks.length === 0" class="p-2 text-gray-500">No tasks found</li>
        </ul>
      </div>
    </div>
  </div>
</template>

