import {
    useEffect,
    useState,
    type FormEvent,
} from 'react';

import type { ProjectMember } from '../../types/project';
import type {
    TaskPriority,
    TaskStatus,
    TaskWriteRequest,
} from '../../types/task';

export type TaskFormValues = {
    name: string;
    description: string;
    assigneeUserId: string;
    status: TaskStatus;
    priority: TaskPriority;
    startDate: string;
    dueDate: string;
};

type TaskFormProps = {
    initialValues?: TaskFormValues;
    members: ProjectMember[];
    isSubmitting: boolean;
    submitErrorMessage: string;
    submitLabel: string;
    onSubmit: (
        request: TaskWriteRequest,
    ) => Promise<void>;
    onCancel: () => void;
};

const defaultValues: TaskFormValues = {
    name: '',
    description: '',
    assigneeUserId: '',
    status: 'not_started',
    priority: 'medium',
    startDate: '',
    dueDate: '',
};

function normalizeNullableValue(
    value: string,
): string | null {
    const normalizedValue = value.trim();

    return normalizedValue || null;
}

export function TaskForm({
    initialValues = defaultValues,
    members,
    isSubmitting,
    submitErrorMessage,
    submitLabel,
    onSubmit,
    onCancel,
}: TaskFormProps) {
    const [values, setValues] =
        useState<TaskFormValues>(
            initialValues,
        );

    const [
        validationMessage,
        setValidationMessage,
    ] = useState('');

    useEffect(() => {
        setValues(initialValues);
    }, [initialValues]);

    async function handleSubmit(
        event: FormEvent<HTMLFormElement>,
    ): Promise<void> {
        event.preventDefault();
        setValidationMessage('');

        const name = values.name.trim();

        const assigneeUserId = Number(
            values.assigneeUserId,
        );

        if (!name) {
            setValidationMessage(
                'タスク名を入力してください。',
            );
            return;
        }

        if (name.length > 255) {
            setValidationMessage(
                'タスク名は255文字以内で入力してください。',
            );
            return;
        }

        if (
            !Number.isSafeInteger(
                assigneeUserId,
            ) ||
            assigneeUserId <= 0
        ) {
            setValidationMessage(
                '担当者を選択してください。',
            );
            return;
        }

        if (
            values.startDate &&
            values.dueDate &&
            values.dueDate <
                values.startDate
        ) {
            setValidationMessage(
                '期限は開始日以降の日付を指定してください。',
            );
            return;
        }

        await onSubmit({
            name,
            description:
                normalizeNullableValue(
                    values.description,
                ),
            assignee_user_id:
                assigneeUserId,
            status: values.status,
            priority: values.priority,
            start_date:
                normalizeNullableValue(
                    values.startDate,
                ),
            due_date:
                normalizeNullableValue(
                    values.dueDate,
                ),
        });
    }

    return (
        <form
            className="task-form"
            onSubmit={(event) => {
                void handleSubmit(event);
            }}
        >
            {(validationMessage ||
                submitErrorMessage) && (
                <div
                    className="task-form__message task-form__message--error"
                    role="alert"
                >
                    {validationMessage ||
                        submitErrorMessage}
                </div>
            )}

            <div className="task-form__fields">
                <label className="task-form__field">
                    <span className="task-form__label">
                        タスク名
                        <span aria-hidden="true">
                            *
                        </span>
                    </span>

                    <input
                        className="task-form__input"
                        type="text"
                        required
                        maxLength={255}
                        value={values.name}
                        onChange={(event) => {
                            setValues(
                                (
                                    currentValues,
                                ) => ({
                                    ...currentValues,
                                    name: event.target
                                        .value,
                                }),
                            );
                        }}
                    />
                </label>

                <label className="task-form__field">
                    <span className="task-form__label">
                        説明
                    </span>

                    <textarea
                        className="task-form__textarea"
                        maxLength={5000}
                        rows={5}
                        value={values.description}
                        onChange={(event) => {
                            setValues(
                                (
                                    currentValues,
                                ) => ({
                                    ...currentValues,
                                    description:
                                        event.target
                                            .value,
                                }),
                            );
                        }}
                    />
                </label>

                <label className="task-form__field">
                    <span className="task-form__label">
                        担当者
                        <span aria-hidden="true">
                            *
                        </span>
                    </span>

                    <select
                        className="task-form__select"
                        required
                        value={
                            values.assigneeUserId
                        }
                        onChange={(event) => {
                            setValues(
                                (
                                    currentValues,
                                ) => ({
                                    ...currentValues,
                                    assigneeUserId:
                                        event.target
                                            .value,
                                }),
                            );
                        }}
                    >
                        <option value="">
                            選択してください
                        </option>

                        {members.map((member) => (
                            <option
                                key={
                                    member.user_id
                                }
                                value={
                                    member.user_id
                                }
                            >
                                {
                                    member.display_name
                                }
                            </option>
                        ))}
                    </select>
                </label>

                <div className="task-form__select-fields">
                    <label className="task-form__field">
                        <span className="task-form__label">
                            ステータス
                            <span aria-hidden="true">
                                *
                            </span>
                        </span>

                        <select
                            className="task-form__select"
                            value={values.status}
                            onChange={(event) => {
                                setValues(
                                    (
                                        currentValues,
                                    ) => ({
                                        ...currentValues,
                                        status: event
                                            .target
                                            .value as TaskStatus,
                                    }),
                                );
                            }}
                        >
                            <option value="not_started">
                                未着手
                            </option>

                            <option value="in_progress">
                                進行中
                            </option>

                            <option value="completed">
                                完了
                            </option>

                            <option value="suspended">
                                保留
                            </option>
                        </select>
                    </label>

                    <label className="task-form__field">
                        <span className="task-form__label">
                            優先度
                            <span aria-hidden="true">
                                *
                            </span>
                        </span>

                        <select
                            className="task-form__select"
                            value={
                                values.priority
                            }
                            onChange={(event) => {
                                setValues(
                                    (
                                        currentValues,
                                    ) => ({
                                        ...currentValues,
                                        priority:
                                            event
                                                .target
                                                .value as TaskPriority,
                                    }),
                                );
                            }}
                        >
                            <option value="high">
                                高
                            </option>

                            <option value="medium">
                                中
                            </option>

                            <option value="low">
                                低
                            </option>
                        </select>
                    </label>
                </div>

                <div className="task-form__date-fields">
                    <label className="task-form__field">
                        <span className="task-form__label">
                            開始日
                        </span>

                        <input
                            className="task-form__input"
                            type="date"
                            value={
                                values.startDate
                            }
                            onChange={(event) => {
                                setValues(
                                    (
                                        currentValues,
                                    ) => ({
                                        ...currentValues,
                                        startDate:
                                            event
                                                .target
                                                .value,
                                    }),
                                );
                            }}
                        />
                    </label>

                    <label className="task-form__field">
                        <span className="task-form__label">
                            期限
                        </span>

                        <input
                            className="task-form__input"
                            type="date"
                            min={
                                values.startDate ||
                                undefined
                            }
                            value={values.dueDate}
                            onChange={(event) => {
                                setValues(
                                    (
                                        currentValues,
                                    ) => ({
                                        ...currentValues,
                                        dueDate:
                                            event
                                                .target
                                                .value,
                                    }),
                                );
                            }}
                        />
                    </label>
                </div>
            </div>

            <footer className="task-form__actions">
                <button
                    className="task-form__button task-form__button--secondary"
                    type="button"
                    disabled={isSubmitting}
                    onClick={onCancel}
                >
                    キャンセル
                </button>

                <button
                    className="task-form__button task-form__button--primary"
                    type="submit"
                    disabled={
                        isSubmitting ||
                        members.length === 0
                    }
                >
                    {isSubmitting
                        ? '処理中です。'
                        : submitLabel}
                </button>
            </footer>
        </form>
    );
}