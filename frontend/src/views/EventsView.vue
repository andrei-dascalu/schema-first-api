<script setup lang="ts">
import { onMounted, ref } from "vue";
import { createEvent, deleteEvent, listEvents } from "../api/apiClient";
import type { Event } from "../api/apiClient";
import { formatDateStub, formatDateTime } from "../lib/date";
import AppInput from "../components/AppInput.vue";
import AppButton from "../components/AppButton.vue";
import AlertBanner from "../components/AlertBanner.vue";
import EmptyState from "../components/EmptyState.vue";

const events = ref<Event[]>([]);
const loading = ref(true);
const error = ref("");

const showCreateForm = ref(false);
const creating = ref(false);
const name = ref("");
const date = ref("");
const location = ref("");

const confirmingDeleteId = ref<string | null>(null);
const deletingId = ref<string | null>(null);

async function loadEvents() {
  loading.value = true;
  error.value = "";
  const { data, error: apiError } = await listEvents();
  loading.value = false;
  if (apiError) {
    error.value = "message" in apiError ? apiError.message : "Couldn't load your events.";
    return;
  }
  events.value = [...data].sort(
    (a, b) => new Date(a.date).getTime() - new Date(b.date).getTime(),
  );
}

async function onCreate() {
  error.value = "";
  creating.value = true;
  const { error: apiError } = await createEvent({
    body: {
      name: name.value,
      date: new Date(date.value).toISOString(),
      location: location.value || undefined,
    },
  });
  creating.value = false;
  if (apiError) {
    error.value = "message" in apiError ? apiError.message : "Couldn't create that event.";
    return;
  }
  name.value = "";
  date.value = "";
  location.value = "";
  showCreateForm.value = false;
  await loadEvents();
}

async function onDelete(eventId: string) {
  deletingId.value = eventId;
  const { error: apiError } = await deleteEvent({ path: { eventId } });
  deletingId.value = null;
  confirmingDeleteId.value = null;
  if (apiError) {
    error.value = "message" in apiError ? apiError.message : "Couldn't delete that event.";
    return;
  }
  events.value = events.value.filter((e) => e.id !== eventId);
}

onMounted(loadEvents);
</script>

<template>
  <div class="page">
    <div class="page-header">
      <h1>Events</h1>
      <AppButton
        :variant="showCreateForm ? 'secondary' : 'primary'"
        @click="showCreateForm = !showCreateForm"
      >
        {{ showCreateForm ? "Cancel" : "New event" }}
      </AppButton>
    </div>

    <Transition name="disclosure">
      <div v-if="showCreateForm" class="disclosure-inner" style="margin-bottom: 20px">
        <form class="panel stack" @submit.prevent="onCreate">
          <div class="form-row cols-3">
            <AppInput v-model="name" label="Name" required />
            <AppInput v-model="date" label="Date & time" type="datetime-local" required />
            <AppInput v-model="location" label="Location" placeholder="Optional" />
          </div>
          <div class="form-actions">
            <AppButton type="submit" :loading="creating">
              {{ creating ? "Creating…" : "Create event" }}
            </AppButton>
          </div>
        </form>
      </div>
    </Transition>

    <AlertBanner v-if="error" :message="error" />

    <p v-if="loading" class="loading-line">
      <span class="spinner" aria-hidden="true" />
      Loading events…
    </p>

    <EmptyState
      v-else-if="events.length === 0"
      title="No events yet"
      message="Create your first event above to start building a guest list."
    />

    <ul v-else class="row-list">
      <li v-for="event in events" :key="event.id" class="row">
        <template v-if="confirmingDeleteId === event.id">
          <p class="row-confirm">
            Delete "{{ event.name }}" and its guest list?
            <AppButton
              variant="danger-solid"
              size="sm"
              :loading="deletingId === event.id"
              @click="onDelete(event.id)"
            >
              Delete
            </AppButton>
            <AppButton variant="ghost" size="sm" @click="confirmingDeleteId = null">
              Cancel
            </AppButton>
          </p>
        </template>
        <template v-else>
          <RouterLink :to="{ name: 'event-detail', params: { eventId: event.id } }" class="row-link">
            <span class="date-stub">
              <span class="day">{{ formatDateStub(event.date).day }}</span>
              <span class="month">{{ formatDateStub(event.date).month }}</span>
            </span>
            <span class="row-body">
              <h2>{{ event.name }}</h2>
              <span class="meta">
                {{ formatDateTime(event.date) }}<template v-if="event.location"> · {{ event.location }}</template>
              </span>
            </span>
          </RouterLink>
          <span class="row-actions">
            <AppButton variant="danger" size="sm" @click="confirmingDeleteId = event.id">
              Delete
            </AppButton>
          </span>
        </template>
      </li>
    </ul>
  </div>
</template>
