import { writable } from 'svelte/store';

export interface Todo {
	id: number;
	title: string;
	completed: boolean;
}

export const syncingTodos = writable(false);
export const todoError = writable<string | null>(null);

export const todos = initializeStore(() => {
	const { set, update, subscribe } = writable(new Array<Todo>());

	return {
		subscribe,
		set,
		addTodo: (todo: Todo) => update((todos) => [...todos, todo]),
		updateTodo: (todo: Todo) =>
			update((todos) =>
				todos.map((existing) => {
					if (existing.id === todo.id) {
						return todo;
					}
					return existing;
				})
			),
		removeTodo: (id: number) => update((todos) => todos.filter((todo) => todo.id !== id))
	};
});

export async function fetchTodos(): Promise<void> {
	return call<Todo[]>('GET', 'todo', todos.set);
}

export async function createTodo(title: string): Promise<void> {
	return call<Todo>('POST', 'todo', todos.addTodo, { title });
}

export async function removeTodo(id: number): Promise<void> {
	return call<void>('DELETE', `todo/${id}`, () => todos.removeTodo(id));
}

export async function updateTodo(todo: Partial<Todo>): Promise<void> {
	return call<Todo>('PUT', `todo/${todo.id}`, todos.updateTodo, todo);
}

async function call<T>(
	method: 'POST' | 'GET' | 'PUT' | 'DELETE',
	path: string,
	handler: (res: T) => void,
	body?: object
): Promise<void> {
	todoError.set(null);
	syncingTodos.set(true);
	return fetch(`http://localhost:8080/v1/api/${path}`, {
		method: method,
		headers: {
			'Content-Type': 'application/json'
		},
		body: body ? JSON.stringify(body) : undefined
	})
		.then(async (response) => {
			if (response.ok) {
				return response.text().then((t) => (t ? JSON.parse(t) : {}));
			}
			console.log('!!', response);
			const { error } = await response.json();
			throw new Error(error);
		})
		.then((response) => handler(response))
		.catch(todoError.set)
		.finally(() => syncingTodos.set(false));
}

/** Helper function to define and export stores in one step */
function initializeStore<T>(func: () => T): T {
	return func();
}
