<script lang="ts">
	import { flip } from 'svelte/animate';
	import { send, receive } from './transition.js';
	import { updateTodo, type Todo, removeTodo } from './todos.js';
	import DustBin from './DustBin.svelte';

	export let todos: Todo[];
</script>

<ul class="todos">
	{#each todos as todo (todo.id)}
		<li
			class="m-2 p-2 shadow"
			class:completed={todo.completed}
			in:receive={{ key: todo.id }}
			out:send={{ key: todo.id }}
			animate:flip={{ duration: 200 }}
		>
			<label class="cursor-pointer">
				<input
					type="checkbox"
					class="cursor-pointer"
					checked={todo.completed}
					on:change={() => updateTodo(todo.id, { completed: !todo.completed })}
				/>

				<span>{todo.title}</span>

				<button aria-label="Remove" on:click={() => removeTodo(todo.id)}>
					<DustBin />
				</button>
			</label>
		</li>
	{/each}
</ul>

<style>
	label {
		width: 100%;
		height: 100%;
		display: flex;
	}

	.completed {
		opacity: 0.5;
	}

	span {
		flex: 1;
	}

	button {
		background-image: url(./remove.svg);
	}
</style>
