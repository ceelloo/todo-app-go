import { Check, Plus, Timer, Trash } from "lucide-react";
import Go from "./assets/go";
import { ThemeSwitcher } from "./components/theme-provider";
import ReactQuery from "./assets/react-query";
import { Badge } from "./components/ui/badge";
import { Input } from "./components/ui/input";
import { Button } from "./components/ui/button";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";

export type Task = {
  _id: number;
  completed: boolean;
  body: string;
};

export const App = () => {
  const { data } = useQuery<Task[]>({
    queryKey: ["tasks"],
    queryFn: async () => {
      const res = await fetch("http://localhost:3000/api/task");
      return await res.json();
    },
  });

  return (
    <div>
      <Header />
      <div className="mx-auto max-w-2xl mt-8">
        <CreateTask />
        <div className="flex flex-col gap-2">
          {data?.map((data, idx) => (
            <TaskCard {...data} key={idx} />
          ))}
        </div>
      </div>
    </div>
  );
};

const Header = () => {
  return (
    <header className="z-50 mx-auto mt-10 max-w-2xl h-16 rounded-xl border">
      <div className="flex gap-4 items-center h-full p-4">
        <div className="flex-1 flex items-center">
          <Go className="size-10" /> <Plus /> <ReactQuery className="size-8" />
        </div>
        <div className="flex-1 flex justify-center">
          <h1 className="text-2xl font-bold">TO DO LIST</h1>
        </div>
        <div className="flex-1 flex justify-end">
          <ThemeSwitcher />
        </div>
      </div>
    </header>
  );
};

export const TaskCard = (task: Task) => {
  const queryClient = useQueryClient();
  const { mutate: updateTask } = useMutation({
    mutationKey: ["updateTask"],
    mutationFn: async () => {
      const res = await fetch(`http://localhost:3000/api/task/${task._id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ completed: !task.completed }),
      });
      return await res.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
    },
  });

  const { mutate: deleteTask } = useMutation({
    mutationKey: ["deleteTask"],
    mutationFn: async () => {
      const res = await fetch(`http://localhost:3000/api/task/${task._id}`, {
        method: "DELETE",
      });
      return await res.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
    },
  });

  return (
    <div
      data-completed={task.completed}
      className="group data-[completed=true]:text-green-600 data-[completed=false]:text-amber-600 flex justify-between w-full p-4 border rounded-xl"
    >
      <h1 className="text-xl font-semibold">{task.body}</h1>
      <div className="flex items-center gap-4">
        <Badge className="group-data-[completed=true]:bg-green-400 group-data-[completed=false]:bg-amber-400">
          {task.completed ? (
            <>
              <Check /> completed
            </>
          ) : (
            <>
              <Timer /> in progress
            </>
          )}
        </Badge>
        {!task.completed && (
          <Button
            size="sm"
            className="size-6 bg-green-400"
            onClick={() => updateTask()}
          >
            <Check />
          </Button>
        )}
        <Button
          size="sm"
          className="size-6 bg-destructive"
          onClick={() => deleteTask()}
        >
          <Trash />
        </Button>
      </div>
    </div>
  );
};

const CreateTask = () => {
  const [taskBody, setTaskBody] = useState("");
  const queryClient = useQueryClient();

  const { mutate, isPending } = useMutation({
    mutationKey: ["createTask"],
    mutationFn: async (body: string) => {
      const res = await fetch("http://localhost:3000/api/task", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ body }),
      });
      if (!res.ok) throw new Error("Failed to create task");
      return await res.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      setTaskBody("");
    },
  });

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (taskBody.trim()) {
      mutate(taskBody);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex gap-4 my-4">
      <Input
        placeholder="Add a task"
        value={taskBody}
        onChange={(e) => setTaskBody(e.target.value)}
      />
      <Button type="submit" disabled={isPending}>
        <Plus className="w-4 h-4" />
      </Button>
    </form>
  );
};
