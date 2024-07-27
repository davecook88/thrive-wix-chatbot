<script lang="ts">
  import { ChevronDown, ChevronUp, Calendar, BookOpen, MessageSquare } from 'lucide-svelte';

  export let emailAddresses: string[];
  export let spanishLevel: 'Beginner' | 'Intermediate' | 'Advanced';
  export let notes: string;
  export let upcomingClasses: { date: string; topic: string }[];

  let isExpanded = false;

  const levelColors = {
    Beginner: 'bg-green-100 text-green-800',
    Intermediate: 'bg-yellow-100 text-yellow-800',
    Advanced: 'bg-red-100 text-red-800',
  } as const;

  $: levelColor = spanishLevel in levelColors ? levelColors[spanishLevel] : '';
</script>

<div class="bg-white shadow-md rounded-lg p-4 mb-4">
  <div class="flex justify-between items-center">
    <div class="flex items-center space-x-4">
      <div>
        <h3 class="text-lg font-semibold">Customer Info</h3>
        <p class="text-sm text-gray-500">{emailAddresses.join(', ')}</p>
      </div>
      <span class="px-2 py-1 rounded-full text-xs font-medium {levelColor}">
        {spanishLevel}
      </span>
    </div>
    <button on:click={() => (isExpanded = !isExpanded)} class="text-gray-500 hover:text-gray-700">
      {#if isExpanded}
        <ChevronUp size={20} />
      {:else}
        <ChevronDown size={20} />
      {/if}
    </button>
  </div>

  {#if isExpanded}
    <div class="mt-4 space-y-4">
      <div>
        <h4 class="text-sm font-medium text-gray-700 flex items-center">
          <MessageSquare size={16} class="mr-2" /> Notes
        </h4>
        <p class="text-sm text-gray-600 mt-1">{notes}</p>
      </div>

      <div>
        <h4 class="text-sm font-medium text-gray-700 flex items-center">
          <Calendar size={16} class="mr-2" /> Upcoming Classes
        </h4>
        <ul class="mt-1 space-y-2">
          {#each upcomingClasses as classInfo}
            <li class="text-sm text-gray-600 flex items-center">
              <BookOpen size={14} class="mr-2" />
              <span>{classInfo.date} - {classInfo.topic}</span>
            </li>
          {/each}
        </ul>
      </div>
    </div>
  {/if}
</div>
