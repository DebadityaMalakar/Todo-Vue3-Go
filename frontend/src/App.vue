<script setup>
import Search from "./components/Search.vue";
import AddTask from "./components/AddTasks.vue";
import TaskTable from "./components/TasksTable.vue";
</script>

<template>
  <head>
    <title>Todo List</title>
  </head>
  <main class="flex flex-col justify-center items-center w-full h-[100dvh] bg-blue-300 space-y-6">
    <Search />
    <AddTask />
    <TaskTable />
  </main>
</template>
