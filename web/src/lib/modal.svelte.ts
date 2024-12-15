import type { ApiClient } from "$lib/api/client";
import type { QueryArtist } from "$lib/types";
import { getContext, setContext } from "svelte";

export type ModalConfirm = {
  type: "modal-confirm";
  title: string;
  description?: string;
  confirmDelete?: boolean;

  onConfirm?: () => void;
};

export type ModalQueryArtist = {
  type: "modal-query-artist";
  title?: string;
  apiClient: ApiClient;

  // eslint-disable-next-line no-unused-vars
  onArtistSelected: (artist: QueryArtist) => void;
};

export type Modal = ModalConfirm | ModalQueryArtist;

type FullModal = {
  id: string;
  data: Modal;
};

export class ModalState {
  modals = $state<FullModal[]>([]);

  pushModal(modal: Modal): string {
    const id = crypto.randomUUID();
    this.modals.push({ id: id, data: modal });

    return id;
  }

  popModal() {
    this.modals.pop();
  }

  removeModal(id: string) {
    this.modals = this.modals.filter((modal) => modal.id !== id);
  }
}

const MODAL_KEY = Symbol("MODAL");

export function setModalState() {
  return setContext(MODAL_KEY, new ModalState());
}

export function getModalState() {
  return getContext<ReturnType<typeof setModalState>>(MODAL_KEY);
}
