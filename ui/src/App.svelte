<script lang="ts">
	import { scale } from 'svelte/transition';
	import Spinner from './lib/Spinner.svelte';
	import TodoList from './lib/TodoList.svelte';
	import { createTodo, fetchTodos, syncingTodos, todoError, todos } from './lib/todos';

	let inputRef: HTMLElement;
	fetchTodos().then(() => inputRef.focus());
</script>

<main>
	<div>
		<div class="flex justify-center space-x-2">
			<h1 class="mb-6">Your Todo Items!</h1>
			<div class="w-8">
				{#if $syncingTodos}
					<Spinner />
				{/if}
			</div>
		</div>
	</div>

	{#if $todoError}
		<p transition:scale class="mb-4 mt-2 bg-red-600 p-2">{$todoError}</p>
	{/if}

	<div class="board m-60 mt-14">
		<input
			bind:this={inputRef}
			class="border"
			placeholder="what needs to be done?"
			on:keydown={(e) => {
				if (e.key !== 'Enter') return;
				createTodo(e.currentTarget.value);
				e.currentTarget.value = '';
			}}
		/>
		<div>
			<h2>Todo</h2>
			<TodoList todos={$todos.filter((t) => !t.completed)} />
		</div>
		<div>
			<h2>Done</h2>
			<TodoList todos={$todos.filter((t) => t.completed)} />
		</div>
	</div>
</main>

<style>
	.board {
		display: grid;
		grid-template-columns: 1fr 1fr;
		grid-column-gap: 1em;
		max-width: 36em;
		margin: 0 auto;
	}

	.board > input {
		font-size: 1.4em;
		grid-column: 1/3;
		padding: 0.5em;
		margin: 0 0 1rem 0;
	}
</style>
