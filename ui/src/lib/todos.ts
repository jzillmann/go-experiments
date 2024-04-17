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
		updateTodo: (id: number, todo: Partial<Omit<Todo, 'id'>>) =>
			update((todos) =>
				todos.map((existing) => {
					if (existing.id === id) {
						return { ...existing, ...todo };
					}
					return existing;
				})
			),
		removeTodo: (id: number) => update((todos) => todos.filter((todo) => todo.id !== id))
	};
});

export async function fetchTodos(): Promise<void> {
	return call<Todo[]>('GET', 'todos', todos.set);
}

export async function createTodo(title: string): Promise<void> {
	return call<Todo>('POST', 'todos', todos.addTodo, { title });
}

export async function removeTodo(id: number): Promise<void> {
	return call<void>('DELETE', `todos/${id}`, () => todos.removeTodo(id));
}

export async function updateTodo(id: number, todo: Partial<Omit<Todo, 'id'>>): Promise<void> {
	return call<Todo>('PATCH', `todos/${id}`, () => todos.updateTodo(id, todo), todo);
}

async function call<T>(
	method: 'POST' | 'GET' | 'PUT' | 'PATCH' | 'DELETE',
	path: string,
	handler: (res: T) => void,
	body?: object
): Promise<void> {
	todoError.set(null);
	syncingTodos.set(true);
	return fetch(`http://localhost:8080/${path}`, {
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

			let responseJson;
			try {
				responseJson = await response.json();
			} catch (e) {
				throw new Error(`${response.status}: ${response.statusText}`);
			}
			// typical huma error
			const { errors } = responseJson; // TODO this needs refinement
			throw new Error(`${errors[0].location}: ${errors[0].message}`);
		})
		.then((response) => handler(response))
		.catch((err) => todoError.set(err.message || err))
		.finally(() => syncingTodos.set(false));
}

/** Helper function to define and export stores in one step */
function initializeStore<T>(func: () => T): T {
	return func();
}
