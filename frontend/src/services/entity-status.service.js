const UPDATED = "updated";
const CREATED = "created";
const DELETED = "deleted";

class EntityStatusService {
   created() {
       return CREATED;
   }

   updated() {
        return UPDATED;
    }

    deleted() {
        return DELETED;
    }

    getMessageByStatus(entityType, status) {
        switch (status) {
            case UPDATED:
                return entityType + " successfully updated!";
            case CREATED:
                return entityType + " successfully created!";
            case DELETED:
                return entityType + " successfully deleted!";
            default:
                return "Something went wrong!";
        }
    }
}

export default new EntityStatusService();