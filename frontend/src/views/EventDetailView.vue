<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import {
  createGuest,
  deleteEvent,
  deleteGuest,
  getEvent,
  listGuests,
  updateEvent,
  updateGuest,
} from "../api/apiClient";
import type { Event, Guest, RsvpStatus } from "../api/apiClient";
import { formatDateTime, toDatetimeLocalValue } from "../lib/date";
import AppInput from "../components/AppInput.vue";
import AppButton from "../components/AppButton.vue";
import AlertBanner from "../components/AlertBanner.vue";
import EmptyState from "../components/EmptyState.vue";

const props = defineProps<{ eventId: string }>();
const router = useRouter();

const event = ref<Event | null>(null);
const eventLoading = ref(true);
const eventError = ref("");

const showEditForm = ref(false);
const savingEdit = ref(false);
const editName = ref("");
const editDate = ref("");
const editLocation = ref("");

const confirmingDeleteEvent = ref(false);
const deletingEvent = ref(false);

const guests = ref<Guest[]>([]);
const guestsLoading = ref(true);
const guestsError = ref("");

const showAddGuestForm = ref(false);
const addingGuest = ref(false);
const guestName = ref("");
const guestEmail = ref("");

const confirmingDeleteGuestId = ref<string | null>(null);
const deletingGuestId = ref<string | null>(null);
const updatingGuestId = ref<string | null>(null);

const statusLabels: Record<RsvpStatus, string> = {
  pending: "Pending",
  accepted: "Accepted",
  declined: "Declined",
};

async function loadEvent() {
  eventLoading.value = true;
  eventError.value = "";
  const { data, error: apiError } = await getEvent({ path: { eventId: props.eventId } });
  eventLoading.value = false;
  if (apiError) {
    eventError.value =
      "message" in apiError ? apiError.message : "Couldn't load this event.";
    return;
  }
  event.value = data;
}

async function loadGuests() {
  guestsLoading.value = true;
  guestsError.value = "";
  const { data, error: apiError } = await listGuests({ path: { eventId: props.eventId } });
  guestsLoading.value = false;
  if (apiError) {
    guestsError.value =
      "message" in apiError ? apiError.message : "Couldn't load the guest list.";
    return;
  }
  guests.value = data;
}

function openEditForm() {
  if (!event.value) return;
  editName.value = event.value.name;
  editDate.value = toDatetimeLocalValue(event.value.date);
  editLocation.value = event.value.location ?? "";
  showEditForm.value = true;
}

async function onSaveEdit() {
  eventError.value = "";
  savingEdit.value = true;
  const { data, error: apiError } = await updateEvent({
    path: { eventId: props.eventId },
    body: {
      name: editName.value,
      date: new Date(editDate.value).toISOString(),
      location: editLocation.value || undefined,
    },
  });
  savingEdit.value = false;
  if (apiError) {
    eventError.value =
      "message" in apiError ? apiError.message : "Couldn't save those changes.";
    return;
  }
  event.value = data;
  showEditForm.value = false;
}

async function onDeleteEvent() {
  deletingEvent.value = true;
  const { error: apiError } = await deleteEvent({ path: { eventId: props.eventId } });
  deletingEvent.value = false;
  if (apiError) {
    eventError.value =
      "message" in apiError ? apiError.message : "Couldn't delete this event.";
    return;
  }
  router.push({ name: "events" });
}

async function onAddGuest() {
  guestsError.value = "";
  addingGuest.value = true;
  const { error: apiError } = await createGuest({
    path: { eventId: props.eventId },
    body: { name: guestName.value, email: guestEmail.value || undefined },
  });
  addingGuest.value = false;
  if (apiError) {
    guestsError.value =
      "message" in apiError ? apiError.message : "Couldn't add that guest.";
    return;
  }
  guestName.value = "";
  guestEmail.value = "";
  showAddGuestForm.value = false;
  await loadGuests();
}

async function onSetRsvp(guest: Guest, status: RsvpStatus) {
  if (guest.rsvpStatus === status) return;
  updatingGuestId.value = guest.id;
  const { data, error: apiError } = await updateGuest({
    path: { eventId: props.eventId, guestId: guest.id },
    body: { name: guest.name, email: guest.email, rsvpStatus: status },
  });
  updatingGuestId.value = null;
  if (apiError) {
    guestsError.value =
      "message" in apiError ? apiError.message : "Couldn't update that guest's RSVP.";
    return;
  }
  const idx = guests.value.findIndex((g) => g.id === guest.id);
  if (idx !== -1) guests.value[idx] = data;
}

async function onDeleteGuest(guestId: string) {
  deletingGuestId.value = guestId;
  const { error: apiError } = await deleteGuest({
    path: { eventId: props.eventId, guestId },
  });
  deletingGuestId.value = null;
  confirmingDeleteGuestId.value = null;
  if (apiError) {
    guestsError.value =
      "message" in apiError ? apiError.message : "Couldn't remove that guest.";
    return;
  }
  guests.value = guests.value.filter((g) => g.id !== guestId);
}

onMounted(() => {
  loadEvent();
  loadGuests();
});
</script>

<template>
  <div class="page">
    <RouterLink to="/events" class="back-link">← Back to events</RouterLink>

    <p v-if="eventLoading" class="loading-line">
      <span class="spinner" aria-hidden="true" />
      Loading event…
    </p>

    <template v-else-if="event">
      <div class="detail-header">
        <div>
          <h1>{{ event.name }}</h1>
          <div class="detail-meta">
            <span class="meta-chip">{{ formatDateTime(event.date) }}</span>
            <span v-if="event.location" class="meta-chip">{{ event.location }}</span>
          </div>
        </div>
        <div class="row-actions">
          <AppButton variant="secondary" size="sm" @click="openEditForm">Edit</AppButton>
          <AppButton variant="danger" size="sm" @click="confirmingDeleteEvent = true">
            Delete
          </AppButton>
        </div>
      </div>

      <Transition name="disclosure">
        <div v-if="showEditForm" class="disclosure-inner" style="margin-bottom: 16px">
          <form class="panel stack" @submit.prevent="onSaveEdit">
            <div class="form-row cols-3">
              <AppInput v-model="editName" label="Name" required />
              <AppInput v-model="editDate" label="Date & time" type="datetime-local" required />
              <AppInput v-model="editLocation" label="Location" placeholder="Optional" />
            </div>
            <div class="form-actions">
              <AppButton type="submit" :loading="savingEdit">
                {{ savingEdit ? "Saving…" : "Save changes" }}
              </AppButton>
              <AppButton variant="ghost" type="button" @click="showEditForm = false">
                Cancel
              </AppButton>
            </div>
          </form>
        </div>
      </Transition>

      <div v-if="confirmingDeleteEvent" class="panel row-confirm" style="margin-bottom: 16px">
        Delete "{{ event.name }}" and its guest list? This can't be undone.
        <AppButton variant="danger-solid" size="sm" :loading="deletingEvent" @click="onDeleteEvent">
          Delete
        </AppButton>
        <AppButton variant="ghost" size="sm" @click="confirmingDeleteEvent = false">
          Cancel
        </AppButton>
      </div>

      <AlertBanner v-if="eventError" :message="eventError" />
    </template>

    <div class="section-heading">
      <h2>Guests</h2>
      <AppButton
        :variant="showAddGuestForm ? 'secondary' : 'primary'"
        size="sm"
        @click="showAddGuestForm = !showAddGuestForm"
      >
        {{ showAddGuestForm ? "Cancel" : "Add guest" }}
      </AppButton>
    </div>

    <Transition name="disclosure">
      <div v-if="showAddGuestForm" class="disclosure-inner" style="margin-bottom: 20px">
        <form class="panel stack" @submit.prevent="onAddGuest">
          <div class="form-row cols-2">
            <AppInput v-model="guestName" label="Name" required />
            <AppInput v-model="guestEmail" label="Email" type="email" placeholder="Optional" />
          </div>
          <div class="form-actions">
            <AppButton type="submit" :loading="addingGuest">
              {{ addingGuest ? "Adding…" : "Add guest" }}
            </AppButton>
          </div>
        </form>
      </div>
    </Transition>

    <AlertBanner v-if="guestsError" :message="guestsError" />

    <p v-if="guestsLoading" class="loading-line">
      <span class="spinner" aria-hidden="true" />
      Loading guests…
    </p>

    <EmptyState
      v-else-if="guests.length === 0"
      title="No guests yet"
      message="Add someone above to start the list."
    />

    <ul v-else class="row-list">
      <li v-for="guest in guests" :key="guest.id" class="row">
        <template v-if="confirmingDeleteGuestId === guest.id">
          <p class="row-confirm">
            Remove {{ guest.name }} from the guest list?
            <AppButton
              variant="danger-solid"
              size="sm"
              :loading="deletingGuestId === guest.id"
              @click="onDeleteGuest(guest.id)"
            >
              Remove
            </AppButton>
            <AppButton variant="ghost" size="sm" @click="confirmingDeleteGuestId = null">
              Cancel
            </AppButton>
          </p>
        </template>
        <template v-else>
          <span class="row-body">
            <h2>{{ guest.name }}</h2>
            <span v-if="guest.email" class="meta">{{ guest.email }}</span>
          </span>
          <span class="row-actions">
            <span class="segmented" role="group" aria-label="RSVP status">
              <button
                v-for="status in (['pending', 'accepted', 'declined'] as RsvpStatus[])"
                :key="status"
                type="button"
                :class="['status-' + status, { 'is-active': guest.rsvpStatus === status }]"
                :disabled="updatingGuestId === guest.id"
                @click="onSetRsvp(guest, status)"
              >
                {{ statusLabels[status] }}
              </button>
            </span>
            <AppButton variant="danger" size="sm" @click="confirmingDeleteGuestId = guest.id">
              Remove
            </AppButton>
          </span>
        </template>
      </li>
    </ul>
  </div>
</template>
